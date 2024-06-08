import { useState, useEffect } from "react";
import { MapContainer, TileLayer, FeatureGroup, Polygon } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import "../index.css"; // Create this file for styling

function App() {
  const [objectsCount, setObjectsCount] = useState(0);
  const [startTime, setStartTime] = useState(null); // Add startTime and endTime to state
  const [endTime, setEndTime] = useState(null);
  const [polygon, setPolygon] = useState([]);
  const [isDrawing, setIsDrawing] = useState(false);

  const position = [21.0278, 105.8342]; // Hanoi coordinates

  // Event handlers for date/time pickers and drawing
  const handleStartTimeChange = (event) =>
    setStartTime(new Date(event.target.value));
  const handleEndTimeChange = (event) =>
    setEndTime(new Date(event.target.value));
  const handlePolygonChange = (polygon) => setPolygon(polygon);

  useEffect(() => {
    const fetchObjectCount = async () => {
      if (startTime && endTime && polygon.length > 0) {
        const startTimestamp = startTime.toISOString();
        const endTimestamp = endTime.toISOString();
        const polygonWKT = `POLYGON((${polygon
          .map((point) => point.join(" "))
          .join(",")}))`; // Convert to WKT

        try {
          const response = await fetch(
            `http://localhost:9090/api/objects/count?start=${startTimestamp}&end=${endTimestamp}&polygon=${encodeURIComponent(
              polygonWKT
            )}`
          );
          if (!response.ok) {
            throw new Error("Network response was not ok.");
          }
          const data = await response.json();
          setObjectsCount(data);
        } catch (err) {
          console.error("Error fetching object count:", err);
        }
      }
    };

    fetchObjectCount();
  }, [startTime, endTime, polygon]);

  return (
    <div className="map">
      <div className="sidebar">
        <div>
          <label htmlFor="start-time">Start Time:</label>
          <input
            type="datetime-local"
            id="start-time"
            onChange={handleStartTimeChange}
          />
        </div>
        <div>
          <label htmlFor="end-time">End Time:</label>
          <input
            type="datetime-local"
            id="end-time"
            onChange={handleEndTimeChange}
          />
        </div>
        <button onClick={() => setIsDrawing(!isDrawing)}>
          {isDrawing ? "Finish Drawing" : "Start Drawing"}
        </button>
        <p>Object Count: {objectsCount}</p>
      </div>

      <MapContainer
        center={position}
        zoom={13}
        style={{ width: "100%", height: "600px" }}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        <FeatureGroup>
          {isDrawing && (
            <Polygon
              positions={polygon}
              color="blue"
              eventHandlers={{
                onEdited: (e) => handlePolygonChange(e.poly.getLatLngs()[0]), // Capture polygon changes
              }}
            />
          )}
        </FeatureGroup>
      </MapContainer>
    </div>
  );
}

export default App;
