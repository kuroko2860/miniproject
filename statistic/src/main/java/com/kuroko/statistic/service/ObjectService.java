package com.kuroko.statistic.service;

import java.time.LocalDateTime;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.kuroko.statistic.entity.ObjectEntity;
import com.kuroko.statistic.repository.ObjectRepository;

@Service
public class ObjectService {

    @Autowired
    private ObjectRepository objectRepository;

    public List<ObjectEntity> getObjectsWithinDistanceAndTimeRange(
            LocalDateTime start,
            LocalDateTime end,
            double longitude,
            double latitude,
            double distance) {
        return objectRepository.findObjectsWithinDistanceAndTimeRange(start, end, longitude, latitude, distance);
    }

    public Long countObjectsWithinPolygonAndTimeRange(LocalDateTime start, LocalDateTime end, String polygonWKT) {
        return objectRepository.countObjectsWithinPolygonAndTimeRange(start, end, polygonWKT);
    }
}