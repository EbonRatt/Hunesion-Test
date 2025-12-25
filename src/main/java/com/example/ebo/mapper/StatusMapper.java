package com.example.ebo.mapper;

import com.example.ebo.model.domain.Status;
import org.apache.ibatis.annotations.Insert;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Update;
import reactor.core.publisher.Mono;

import java.util.UUID;

@Mapper
public interface StatusMapper {

    @Insert("""
              insert into status(player_id, level, exp, hp)
              values(#{status.playerId}, #{status.level}, #{status.exp}, #{status.hp})
              on conflict (player_id) do nothing
            """)
    Mono<Integer> insertIfAbsent(@Param("status") Status status);

    @Update("""
              update status
              set level = level + 1,
                  updated_at = now()
              where player_id = #{playerId}
            """)
    Mono<Integer> levelUp(UUID playerId);

}
