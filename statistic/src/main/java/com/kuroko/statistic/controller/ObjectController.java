package com.kuroko.statistic.controller;

import org.locationtech.jts.geom.GeometryFactory;
import org.locationtech.jts.geom.Polygon;
import org.locationtech.jts.geom.Coordinate;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import com.kuroko.statistic.entity.ObjectEntity;
import com.kuroko.statistic.service.ObjectService;

import java.time.Instant;
import java.util.List;

@CrossOrigin(origins = "*", maxAge = 3600)
@RestController
@RequestMapping("/api/objects")
public class ObjectController {

    @Autowired
    private ObjectService objectService;

    @GetMapping
    public List<ObjectEntity> getObjects(
            @RequestParam double[][] coordinates,
            @RequestParam String startTime,
            @RequestParam String endTime,
            @RequestParam(required = false) String type,
            @RequestParam(required = false) String color) {

        GeometryFactory geometryFactory = new GeometryFactory();
        Coordinate[] coords = new Coordinate[coordinates.length];
        for (int i = 0; i < coordinates.length; i++) {
            coords[i] = new Coordinate(coordinates[i][0], coordinates[i][1]);
        }
        Polygon polygon = geometryFactory.createPolygon(coords);

        Instant start = Instant.parse(startTime);
        Instant end = Instant.parse(endTime);

        return objectService.getObjects(polygon, start, end, type, color);
    }
}
