package com.kuroko.statistic.service;

import org.locationtech.jts.geom.Polygon;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.kuroko.statistic.entity.ObjectEntity;
import com.kuroko.statistic.repository.ObjectRepository;

import java.time.Instant;
import java.util.List;

@Service
public class ObjectService {

    @Autowired
    private ObjectRepository objectRepository;

    public List<ObjectEntity> getObjects(Polygon polygon, Instant startTime, Instant endTime, String type,
            String color) {
        return objectRepository.findByGeoAndTimeRangeAndAttributes(polygon, startTime, endTime, type, color);
    }
}
