.PHONY: deploy delete

deploy:
	gcloud functions deploy subscribe_$(environment)_$(commit) --entry-point Subscribe --runtime go113 --trigger-http --memory 128MB --region europe-west1

delete:
	gcloud functions delete subscribe_$(environment)_$(commit) --region europe-west1

deploy_preprod:
	gcloud functions deploy subscribe_preprod --entry-point Subscribe --runtime go113 --trigger-http --memory 128MB --region europe-west1

test:
	go test .

integration: # TODO add integration tests
	go test .

smoke: # TODO add smoke tests
	go test .