package com.example.ebo.service;

import com.example.ebo.model.dto.request.PlayerCreateRequest;
import com.example.ebo.model.dto.response.PlayerResponse;
import com.example.ebo.model.dto.response.PlayerStatusRowResponse;
import reactor.core.publisher.Mono;

import java.util.UUID;

public interface PlayerService {
    Mono<PlayerResponse> createPlayer(PlayerCreateRequest request);
    Mono<PlayerStatusRowResponse> searchPlayerByUsername(String username);

    Mono<PlayerStatusRowResponse> getPlayerStatus(UUID playerId);

    Mono<PlayerStatusRowResponse> levelUpStatus(UUID playerId);
}
