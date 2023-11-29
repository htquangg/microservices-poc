package am

func MessagePublisherWithMiddleware(publisher MessagePublisher, mws ...MessagePublisherMiddleware) MessagePublisher {
	return applyMiddleware(publisher, mws...)
}

func MessageHandlerWithMiddleware(subscriber MessageHandler, mws ...MessageHandlerMiddleware) MessageHandler {
	return applyMiddleware(subscriber, mws...)
}

func applyMiddleware[T any, M func(T) T](target T, mws ...M) T {
	h := target

	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}
