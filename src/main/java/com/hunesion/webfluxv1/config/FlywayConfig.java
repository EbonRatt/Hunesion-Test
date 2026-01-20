package com.hunesion.webfluxv1.config;

import org.flywaydb.core.Flyway;
import org.flywaydb.core.api.MigrationInfo;
import org.flywaydb.core.api.MigrationInfoService;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.jdbc.datasource.DriverManagerDataSource;

import javax.sql.DataSource;

@Configuration
public class FlywayConfig {

    @Value("${spring.datasource.url}")
    private String url;

    @Value("${spring.datasource.username}")
    private String username;

    @Value("${spring.datasource.password}")
    private String password;

    @Bean
    public DataSource dataSource() {
        DriverManagerDataSource dataSource = new DriverManagerDataSource();
        dataSource.setDriverClassName("org.mariadb.jdbc.Driver");
        dataSource.setUrl(url);
        dataSource.setUsername(username);
        dataSource.setPassword(password);
        return dataSource;
    }

    @Bean(initMethod = "migrate")
    public Flyway flyway(DataSource dataSource) {
        Flyway flyway = Flyway.configure()
                .dataSource(dataSource)
                .baselineOnMigrate(true)
                .locations("classpath:db/migration")
                .cleanDisabled(false)  // Enable clean operation
                .outOfOrder(false)
                .baselineVersion("1")
                .baselineDescription("Initial schema")
                .validateOnMigrate(true)
                .load();

        // Clean and repair if there are any failed migrations
        if (hasFailedMigration(flyway.info())) {
            flyway.repair();
            flyway.migrate();
        }

        return flyway;
    }

    private boolean hasFailedMigration(MigrationInfoService info) {
        for (MigrationInfo i : info.all()) {
            if (i.getState().isFailed()) {
                return true;
            }
        }
        return false;
    }
}
