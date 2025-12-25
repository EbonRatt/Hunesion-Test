package com.example.ebo.SSE;

import com.example.ebo.model.dto.response.PlayerStatusRowResponse;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Sinks;

import java.util.UUID;

@Component
public class PlayerStatusBus {
    private final Sinks.Many<PlayerStatusRowResponse> sink =
            Sinks.many().multicast().onBackpressureBuffer();

    public void publish(PlayerStatusRowResponse e) {
        sink.tryEmitNext(e);
    }

    public Flux<PlayerStatusRowResponse> streamByPlayer(UUID playerId) {
        return sink.asFlux().filter(e -> e != null && e.getId() != null && e.getId().equals(playerId));
    }
}
