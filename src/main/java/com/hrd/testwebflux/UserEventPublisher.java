package com.hrd.testwebflux;

import org.springframework.stereotype.Component;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Sinks;

@Component
public class UserEventPublisher {

    private final Sinks.Many<UserEvent> sink =
            Sinks.many().multicast().onBackpressureBuffer();


    public void publish(UserEvent event) {
        sink.tryEmitNext(event);
    }


    public Flux<UserEvent> getEvents() {
        return sink.asFlux();
    }

}
