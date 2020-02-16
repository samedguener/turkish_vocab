.PHONY: deploy delete

deploy:
	gcloud functions deploy hello --entry-point Hello --runtime go113 --trigger-http --memory 128MB --region europe-west1

delete:
	gcloud functions delete hello --entry-point Hello --runtime go113 --trigger-http --region europe-west1