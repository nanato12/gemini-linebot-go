.PHONY: up
up:
	rm -rf gemini-line-bot-gcf.zip
	zip -r gemini-line-bot-gcf.zip * -x ".*" "Makefile" "*.zip"
	gcloud storage cp gemini-line-bot-gcf.zip gs://nanato12-sandbox-gcf/gemini-line-bot-gcf.zip
