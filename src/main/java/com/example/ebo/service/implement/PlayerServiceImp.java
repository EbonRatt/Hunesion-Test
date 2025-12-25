package com.example.ebo.service.implement;

import com.example.ebo.SSE.PlayerStatusBus;
import com.example.ebo.mapper.PlayerMapper;
import com.example.ebo.mapper.StatusMapper;
import com.example.ebo.model.domain.Player;
import com.example.ebo.model.domain.Status;
import com.example.ebo.model.dto.request.PlayerCreateRequest;
import com.example.ebo.model.dto.response.PlayerResponse;
import com.example.ebo.model.dto.response.PlayerStatusRowResponse;
import com.example.ebo.service.PlayerService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.reactive.TransactionalOperator;
import reactor.core.publisher.Mono;

import java.util.UUID;


@Service
@Slf4j
@RequiredArgsConstructor
public class PlayerServiceImp implements PlayerService {
    private final PlayerMapper playerMapper;
    private final StatusMapper statusMapper;
    private final PlayerStatusBus bus;
    private final TransactionalOperator tx;


    @Override
    public Mono<PlayerResponse> createPlayer(PlayerCreateRequest request) {

        return tx.execute(txSpec ->
                playerMapper.insert(Player.builder()
                                .username(request.getUsername())
                                .build())
                        .flatMap(player -> {
                            log.info("Player created: {}", player);
                            return statusMapper.insertIfAbsent(Status.builder()
                                    .playerId(player.getId())
                                    .level(1)
                                    .exp(0L)
                                    .hp(100)
                                    .build())
                                    .then(playerMapper.findRowById(player.getId()))
                                    .doOnNext(bus::publish)
                                    .then(playerMapper.findById(player.getId()));
                        })
                        .map(player -> PlayerResponse.builder()
                                .id(player.getId())
                                .username(player.getUsername())
                                .createdAt(player.getCreatedAt())
                                .build())
        ).single();
    }

    @Override
    public Mono<PlayerStatusRowResponse> searchPlayerByUsername(String username) {
        return tx.execute(txSpec ->
                        playerMapper.findByUsername(username)
                ).singleOrEmpty()
                .switchIfEmpty(Mono.error(new RuntimeException("Player not found")));
    }

    @Override
    public Mono<PlayerStatusRowResponse> getPlayerStatus(UUID playerId) {
        return tx.execute(txSpec ->
                        playerMapper.findRowById(playerId)
                ).singleOrEmpty()
                .switchIfEmpty(Mono.error(new RuntimeException("Player not found")));
    }

    @Override
    public Mono<PlayerStatusRowResponse> levelUpStatus(UUID playerId) {
        return tx.execute(txSpec ->
                        statusMapper.levelUp(playerId)
                                .flatMap(rows -> {
                                    if (rows == null || rows == 0) {
                                        return Mono.error(new RuntimeException("Player status not found"));
                                    }
                                    return playerMapper.findRowById(playerId);
                                })
                                .doOnNext(bus::publish)
                ).single();
    }
}
