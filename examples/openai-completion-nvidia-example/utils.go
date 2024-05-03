package main

import "os"

func SetUp() {
	os.Setenv("LANGCHAIN_TRACING", "true")
	os.Setenv("LANGSMITH_API_KEY", "ls__9e0570672e5241cea39d7200c4f422ea")
	os.Setenv("LANGCHAIN_PROJECT_NAME", "openai-completion-nvidia-example")
	os.Setenv("JINA_API_KEY", "jina_eda6a00a90ac48daabac72b6d3ba5e3d7Dl_rdQSZwyZ04aRdkcVYIzOjtd7")
	os.Setenv("MARITACA_KEY", "116948195544263252464$f410151375330261")
	os.Setenv("OPENAI_NVAPI_KEY", "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")
}
