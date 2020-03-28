.PHONY: deploy delete

deploy:
	gcloud functions deploy subscriptions$(commit) --entry-point Subscribe --runtime go113 --trigger-http --memory 128MB --region $(region) --service-account $(serviceaccount) --impersonate-service-account $(serviceaccount)

delete:
	gcloud functions delete subscriptions$(commit) --region $(region)

test:
	echo ${GOOGLE_APPLICATION_CREDENTIALS}
	go test .

integration: # TODO add integration tests

smoke: # TODO add smoke tests
