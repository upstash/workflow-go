package workflow

const sdkVersion = "1"

const (
	productionUrl = "https://qstash.upstash.io/"
)

const (
	forwardPrefix = "Upstash-Forward-"

	authorizationHeader = "Authorization"
	contentTypeHeader   = "Content-Type"

	initHeader             = "Upstash-Workflow-Init"
	runIdHeader            = "Upstash-Workflow-Runid"
	urlHeader              = "Upstash-Workflow-Url"
	sdkVersionHeader       = "Upstash-Workflow-Sdk-Version"
	featureSetHeader       = "Upstash-Feature-Set"
	flowControlKeyHeader   = "Upstash-Flow-Control-Key"
	flowControlValueHeader = "Upstash-Flow-Control-Value"
	retriesHeader          = "Upstash-Retries"

	featureLazyFetch   = "LazyFetch"
	featureInitialBody = "InitialBody"
)
