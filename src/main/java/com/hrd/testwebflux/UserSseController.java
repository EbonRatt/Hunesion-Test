package com.hrd.testwebflux;

import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.codec.ServerSentEvent;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Flux;

@RestController
@RequestMapping("/sse")
@RequiredArgsConstructor
public class UserSseController {

    private final UserEventPublisher publisher;

    @CrossOrigin(origins = "http://localhost:3000")
    @GetMapping(value = "/users", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Flux<ServerSentEvent<UserEvent>> stream() {
        return publisher.getEvents()
                .map(data -> ServerSentEvent.builder(data).build());
//                .mergeWith(Flux.never()); // keep the connection open
    }

}
