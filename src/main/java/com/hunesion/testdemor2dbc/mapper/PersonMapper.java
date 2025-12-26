package com.hunesion.testdemor2dbc.mapper;

import com.hunesion.testdemor2dbc.model.dto.response.PersonResourceResponse;
import com.hunesion.testdemor2dbc.model.entity.Person;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface PersonMapper {

    PersonResourceResponse toResource(Person person);

}
