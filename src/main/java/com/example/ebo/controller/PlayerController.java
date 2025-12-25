package com.example.ebo.controller;

import com.example.ebo.SSE.PlayerStatusBus;
import com.example.ebo.model.dto.request.PlayerCreateRequest;
import com.example.ebo.model.dto.response.PlayerResponse;
import com.example.ebo.model.dto.response.PlayerStatusRowResponse;
import com.example.ebo.service.PlayerService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.http.codec.ServerSentEvent;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.UUID;

@RestController
@RequestMapping("/players")
@RequiredArgsConstructor
public class PlayerController {
    private final PlayerService playerService;
    private final PlayerStatusBus bus;

    @PostMapping
    public Mono<ResponseEntity<PlayerResponse>> createPlayer(@RequestBody PlayerCreateRequest request) {
        return playerService.createPlayer(request)
                .map(ResponseEntity::ok);
    }

    @GetMapping("/{username}")
    public Mono<ResponseEntity<PlayerStatusRowResponse>> searchPlayerByUsername(@PathVariable String username) {
        return playerService.searchPlayerByUsername(username)
                .map(ResponseEntity::ok);
    }

    @GetMapping(value = "/{playerId}/status/stream", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Flux<ServerSentEvent<PlayerStatusRowResponse>> streamPlayerStatus(@PathVariable UUID playerId) {
        Flux<ServerSentEvent<PlayerStatusRowResponse>> initial = playerService.getPlayerStatus(playerId)
                .flux()
                .map(e -> ServerSentEvent.builder(e)
                        .event("player-status")
                        .id(e.getId() == null ? null : e.getId().toString())
                        .build());

        Flux<ServerSentEvent<PlayerStatusRowResponse>> updates = bus.streamByPlayer(playerId)
                .map(e -> ServerSentEvent.builder(e)
                        .event("player-status")
                        .id(e.getId() == null ? null : e.getId().toString())
                        .build());

        return Flux.concat(initial, updates);
    }

    @PostMapping("/{playerId}/status/level-up")
    public Mono<ResponseEntity<PlayerStatusRowResponse>> levelUp(@PathVariable UUID playerId) {
        return playerService.levelUpStatus(playerId)
                .map(ResponseEntity::ok);
    }

}
