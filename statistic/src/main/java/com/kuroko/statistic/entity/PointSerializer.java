package com.kuroko.statistic.entity;

import java.io.*;

import org.locationtech.jts.geom.Point;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.*;

// Custom Serializer for Point
class PointSerializer extends JsonSerializer<Point> {
    @Override
    public void serialize(Point value, JsonGenerator gen, SerializerProvider serializers) throws IOException {
        gen.writeStartObject();
        gen.writeNumberField("longitude", value.getX());
        gen.writeNumberField("latitude", value.getY());
        gen.writeEndObject();
    }
}