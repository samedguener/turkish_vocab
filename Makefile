.PHONY: deploy delete

deploy:
	gcloud functions deploy subscribe_$(environment) --entry-point Subscribe --runtime go113 --trigger-http --memory 128MB --region europe-west1

delete:
	gcloud functions delete subscribe_$(environment) --region europe-west1

test:
	go test .

integration: # TODO add integration tests
	go test .

smoke: # TODO add smoke tests
	go test .