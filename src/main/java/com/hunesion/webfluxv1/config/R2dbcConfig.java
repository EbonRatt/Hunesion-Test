package com.hunesion.webfluxv1.config;

import com.hunesion.webfluxv1.utils.UuidUtil;
import io.r2dbc.spi.ConnectionFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.convert.converter.Converter;
import org.springframework.data.convert.ReadingConverter;
import org.springframework.data.convert.WritingConverter;
import org.springframework.data.r2dbc.config.EnableR2dbcAuditing;
import org.springframework.data.r2dbc.convert.R2dbcCustomConversions;
import org.springframework.data.r2dbc.dialect.DialectResolver;
import org.springframework.data.r2dbc.dialect.R2dbcDialect;
import org.springframework.r2dbc.connection.R2dbcTransactionManager;
import org.springframework.transaction.ReactiveTransactionManager;
import org.springframework.transaction.reactive.TransactionalOperator;

import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Configuration
@EnableR2dbcAuditing
public class R2dbcConfig {

    @Bean
    public R2dbcCustomConversions r2dbcCustomConversions(ConnectionFactory connectionFactory) {
        R2dbcDialect dialect = DialectResolver.getDialect(connectionFactory);
        List<Object> converters = new ArrayList<>(dialect.getConverters());

        // Add custom converters only once
        converters.add(new UuidToByteArrayConverter());
        converters.add(new ByteArrayToUuidConverter());

        return R2dbcCustomConversions.of(dialect, converters);
    }

    @ReadingConverter
    public static class ByteArrayToUuidConverter implements Converter<byte[], UUID> {
        @Override
        public UUID convert(byte[] source) {
            return UuidUtil.fromBytes(source);
        }
    }

    @WritingConverter
    public static class UuidToByteArrayConverter implements Converter<UUID, byte[]> {
        @Override
        public byte[] convert(UUID source) {
            return UuidUtil.toBytes(source);
        }
    }
}