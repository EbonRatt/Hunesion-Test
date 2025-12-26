package com.hunesion.testdemor2dbc.model.dto.response;

import com.hunesion.testdemor2dbc.enums.ItemStatus;
import com.hunesion.testdemor2dbc.model.entity.Tag;
import lombok.Data;
import lombok.experimental.Accessors;

import java.time.LocalDateTime;
import java.util.List;

@Data
@Accessors(chain = true)
public class ItemResourceResponse {

    private Long id;
    private Long version;

    private String description;
    private ItemStatus status;

    private PersonResourceResponse assignee;
    private List<Tag> tags;

    private LocalDateTime createdDate;
    private LocalDateTime lastModifiedDate;

}
