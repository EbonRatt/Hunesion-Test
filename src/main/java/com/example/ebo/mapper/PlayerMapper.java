package com.example.ebo.mapper;

import com.example.ebo.model.domain.Player;
import com.example.ebo.model.dto.response.PlayerStatusRowResponse;
import org.apache.ibatis.annotations.*;
import reactor.core.publisher.Mono;

import java.util.UUID;

@Mapper
public interface PlayerMapper {

    @Select(
            """
                SELECT * FROM player WHERE id = #{id}
            """
    )
    @Results(value = {
            @Result(property = "createdAt", column = "created_at")
    })
    Mono<Player> findById(UUID id);

    @Select(
            """
                INSERT INTO player (username)
                VALUES (#{username})
                RETURNING *
            """
    )
    @Results(value = {
            @Result(property = "createdAt", column = "created_at")
    })
    Mono<Player> insert(Player player);


    @Select(
            """
                SELECT
                    p.id,
                    p.username,
                    p.created_at,
                    s.level AS status_level,
                    s.exp AS status_exp,
                    s.hp AS status_hp,
                    s.updated_at AS status_updated_at
                FROM player p
                LEFT JOIN status s ON p.id = s.player_id
                WHERE p.username = #{username}
            """
    )
    @Results(value = {
            @Result(property = "createdAt", column = "created_at"),
            @Result(property = "statusLevel", column = "status_level"),
            @Result(property = "statusExp", column = "status_exp"),
            @Result(property = "statusHp", column = "status_hp"),
            @Result(property = "statusUpdatedAt", column = "status_updated_at")
    })
    Mono<PlayerStatusRowResponse> findByUsername(String username);


    @Select(
            """
                SELECT
                    p.id,
                    p.username,
                    p.created_at,
                    s.level AS status_level,
                    s.exp AS status_exp,
                    s.hp AS status_hp,
                    s.updated_at AS status_updated_at
                FROM player p
                LEFT JOIN status s ON p.id = s.player_id
                WHERE p.id = #{playerId}
            """
    )
    @Results(value = {
            @Result(property = "createdAt", column = "created_at"),
            @Result(property = "statusLevel", column = "status_level"),
            @Result(property = "statusExp", column = "status_exp"),
            @Result(property = "statusHp", column = "status_hp"),
            @Result(property = "statusUpdatedAt", column = "status_updated_at")
    })
    Mono<PlayerStatusRowResponse> findRowById(UUID playerId);
}
