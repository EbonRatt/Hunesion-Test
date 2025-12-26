package com.hunesion.testdemor2dbc.model.entity;

import com.hunesion.testdemor2dbc.enums.ItemStatus;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import lombok.ToString;
import lombok.experimental.Accessors;
import org.springframework.data.annotation.*;
import org.springframework.data.relational.core.mapping.Table;

import java.time.LocalDateTime;
import java.util.List;

@Table
@Getter
@Setter
@Accessors(chain = true)
@NoArgsConstructor
@ToString
public class Item {

    public Item(Long id, Long version) {
        this.id = id;
        this.version = version;
    }

    @Id
    private Long id;

    @Version
    private Long version;

    @Size(max=4000)
    @NotBlank
    private String description;

    @NotNull
    private ItemStatus status = ItemStatus.TODO;

    private Long assigneeId;

    @Transient
    private Person assignee;

    @Transient
    private List<Tag> tags;

    @CreatedDate
    private LocalDateTime createdDate;

    @LastModifiedDate
    private LocalDateTime lastModifiedDate;

}
