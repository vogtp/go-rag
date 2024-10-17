package vecdb

type Option func(*VecDB)

func WithChromaAddress(addr string) Option {
	return func(vd *VecDB) {
		vd.slog = vd.slog.With("chroma", addr)
		vd.chromaAddr = addr
	}
}

func WithOllamaAddress(addr string) Option {
	return func(vd *VecDB) {
		vd.slog = vd.slog.With("ollama", addr)
		vd.ollamaAddr = addr
	}
}

func WithEmbeddingsModel(model string) Option {
	return func(vd *VecDB) {
		vd.slog = vd.slog.With("model", model)
		vd.embeddingsModel = model
	}
}
