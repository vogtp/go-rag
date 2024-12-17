package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/amikos-tech/chroma-go/types"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
	"github.com/vogtp/rag/pkg/cfg"
)

func chromaVecDBOwn(ctx context.Context, index string) error {

	loader := documentloaders.NewNotionDirectory("/home/vogtp/go/src/gitlab-int.its.unibas.ch/vogtp/chatbot/modelfiles/vogtp/")

	slog.Info("Searching vecDB", "index", index)
	model := viper.GetString(cfg.ModelDefault)
	llm, err := getOllamaClient(model)
	if err != nil {
		return fmt.Errorf("cannot load embedding model %s: %w", model, err)
	}

	_, e, err := getEmbedding("mxbai-embed-large")
	if err != nil {
		return err
	}
	store, err := chroma.New(
		chroma.WithChromaURL(chromeURL),
		chroma.WithNameSpace(index),
		chroma.WithEmbedder(e),
		chroma.WithDistanceFunction(types.COSINE),
	)
	if err != nil {
		return fmt.Errorf("cannot create chroma client: %w", err)
	}
	docs, err := loader.LoadAndSplit(ctx, textsplitter.NewRecursiveCharacter())
	if err != nil {
		return err
	}
	if len(docs) < 1 {
		return fmt.Errorf("No documents")
	}
	idfuc := vectorstores.WithIDGenerater(func(ctx context.Context, doc schema.Document) string {
		//fmt.Printf("%+v\n", doc.Metadata)
		return fmt.Sprint(doc.Metadata["source"])
	})
	if _, err := store.AddDocuments(ctx, docs, idfuc); err != nil {
		return fmt.Errorf("cannot add docs: %w", err)
	}
	questions := []string{
		"how can I get a guest account",
		"lost my password",
		"cannot connect",
	}
	for _, question := range questions {
		docs, err := store.SimilaritySearch(ctx, question, 3, vectorstores.WithScoreThreshold(0.3))

		if err != nil {
			return fmt.Errorf("cannot search the docs: %w", err)
		}
		fmt.Printf("\n**************\nQuestion: %s\nDocs: %v\n", question, len(docs))
		// for i, d := range docs {
		// 	fmt.Printf("Doc %v score: %v -> %v %v\n", i, d.Score, d.PageContent, d.Metadata)
		// }

		mem := memory.NewConversationBuffer()
		rec := vectorstores.ToRetriever(
			store,
			5,
			// vectorstores.WithNameSpace(index),
			vectorstores.WithScoreThreshold(0.2),
		)
		c := chains.NewConversationalRetrievalQAFromLLM(llm, rec, mem)
		// llmChain := chains.NewConversation(llm, mem)
		// // condenseChain:=chains.NewStuffDocuments()
		// cov := chains.NewConversationalRetrievalQA(llmChain, nil, rec, mem)
		// endChain := chains.NewRefineDocuments(&conversationChain, vecDBChain)

		// ############################################### STREAMING FUNC
		streamingFuncOpt := chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			_, err := fmt.Print(string(chunk))
			return err
		})

		fmt.Printf("\nRes: %v\n", "ConversationalRetrieval")
		_, err = chains.Run(ctx, c, question, streamingFuncOpt)
		if err != nil {
			return fmt.Errorf("cannot run chain: %w", err)
		}


		fmt.Printf("\nRes2: %v\n", "Chain with prompt")
		llmChain := chains.NewLLMChain(llm, prompts.NewPromptTemplate(" ddd", nil))
		_, err = chains.Run(
			ctx,
			chains.NewRetrievalQAFromLLM(
				llmChain.LLM,
				vectorstores.ToRetriever(
					store,
					5,
					//vectorstores.WithNameSpace(index),
					//vectorstores.WithScoreThreshold(0.8),
				),
			),
			question, streamingFuncOpt,
			//"City with a population of more than 5",
		)
		if err != nil {
			return fmt.Errorf("cannot run chain: %w", err)
		}

		stuffQAChain := chains.LoadStuffQA(llm)

		fmt.Printf("\nRes3: %v\n", "QA stuff")
		_, err = chains.Call(context.Background(), stuffQAChain, map[string]any{
			"input_documents": docs,
			"question":        question,
		}, streamingFuncOpt)
		if err != nil {
			return err
		}
	}
	return nil
}
