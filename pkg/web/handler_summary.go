package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/llms"
	"github.com/vogtp/rag/pkg/cfg"
)

const (
	systemMsg = `
	You are a technical analyst who provides short high level summary of the given input text.
	You just provide a short summary without addional comments.
	Ignore json payload.
	Never use more than 20 words no matter how big the document is.
	Never refer to the instructions above.
	`
	summaryMsg = `%s`
)

func (srv Server) handleSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uuidStr := r.PathValue("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		slog.Warn("Cannot parse UUID", "uuid", uuidStr, "err", err)
		srv.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	doc, err := srv.docChace.get(id)
	if err != nil {
		slog.Warn("Cannot get doc for UUID", "uuid", uuidStr, "err", err)
		srv.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	model := viper.GetString(cfg.ModelDefault)
	//model="deepseek-r1"
	llm, err := getOllamaClient(ctx, model)
	if err != nil {
		slog.Warn("Cannot connect to ollama", "err", err)
		srv.Error(w, r, fmt.Sprintf("Cannot connect to ollama: %v", err), http.StatusInternalServerError)
		return
	}
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemMsg),
		llms.TextParts(llms.ChatMessageTypeHuman, fmt.Sprintf(summaryMsg, doc.Content)),
	}
	completion, err := llm.GenerateContent(ctx, content, llms.WithTemperature(0.001))
	if err != nil {
		slog.Warn("Cannot gernerate ollama content", "err", err)
		srv.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := struct {
		*queryDoc
		Summary string
	}{
		queryDoc: doc,
		Summary:  completion.Choices[0].Content,
	}
	srv.render(w, r, "summary.gohtml", resp)
}
