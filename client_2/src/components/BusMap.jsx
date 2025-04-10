// import React, { useState, useEffect } from 'react';
// import { MapContainer, TileLayer, Marker, Popup, Polyline } from 'react-leaflet';
// import { useWebSocket } from '../context/WebSocketContext';
// import L from 'leaflet';

// // Fix for marker icons in Leaflet with React
// delete L.Icon.Default.prototype._getIconUrl;
// L.Icon.Default.mergeOptions({
//     iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
//     iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
//     shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
// });

// // Custom marker icons
// const busIcon = new L.DivIcon({
//     className: 'bus-marker',
//     html: '<i class="bi bi-bus-front"></i>',
//     iconSize: [40, 40],
//     iconAnchor: [20, 20]
// });

// const stopIcon = new L.DivIcon({
//     className: 'stop-marker',
//     html: '<i class="bi bi-geo-alt-fill"></i>',
//     iconSize: [20, 20],
//     iconAnchor: [10, 10]
// });

// const BusMap = ({ buses = [], stops = [], route = null, center = [37.7749, -122.4194], zoom = 13 }) => {
//     const { buses: liveBuses } = useWebSocket();
//     const [mapBuses, setMapBuses] = useState(buses);

//     useEffect(() => {
//         if (liveBuses.length > 0) {
//             setMapBuses(prevBuses => {
//                 // Update buses with live data where possible
//                 const updatedBuses = [...prevBuses].map(bus => {
//                     const liveBus = liveBuses.find(lb => lb.id === bus.id);
//                     return liveBus ? { ...bus, ...liveBus } : bus;
//                 });
//                 return updatedBuses;
//             });
//         }
//     }, [liveBuses]);

//     // If no buses provided, use live buses
//     useEffect(() => {
//         if (buses.length === 0 && liveBuses.length > 0) {
//             setMapBuses(liveBuses);
//         } else if (buses.length > 0) {
//             setMapBuses(buses);
//         }
//     }, [buses, liveBuses]);

//     return (
//         <div className="map-container">
//             <MapContainer center={center} zoom={zoom} style={{ height: '100%', width: '100%' }}>
//                 <TileLayer
//                     attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
//                     url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
//                 />

//                 {/* Draw bus route if provided */}
//                 {route && route.path && route.path.length > 0 && (
//                     <Polyline
//                         positions={route.path.map(point => [point.latitude, point.longitude])}
//                         color={route.color || '#3388ff'}
//                         weight={5}
//                         opacity={0.7}
//                     />
//                 )}

//                 {/* Display bus stops */}
//                 {stops.map(stop => (
//                     <Marker
//                         key={`stop-${stop.id}`}
//                         position={[stop.latitude, stop.longitude]}
//                         icon={stopIcon}
//                     >
//                         <Popup>
//                             <div>
//                                 <h6>{stop.name}</h6>
//                                 <p>Stop #{stop.order}</p>
//                             </div>
//                         </Popup>
//                     </Marker>
//                 ))}

//                 {/* Display buses */}
//                 {mapBuses.map(bus => (
//                     <Marker
//                         key={`bus-${bus.id}`}
//                         position={[bus.latitude, bus.longitude]}
//                         icon={busIcon}
//                     >
//                         <Popup>
//                             <div>
//                                 <h6>Bus #{bus.busNumber}</h6>
//                                 <p>Route: {bus.route?.name || 'Unknown'}</p>
//                                 <p>Status: <span className={`badge status-${bus.status?.toLowerCase().replace(/\s/g, '-')}`}>
//                                     {bus.status || 'Unknown'}
//                                 </span></p>
//                                 <p>Last Updated: {bus.lastUpdated ? new Date(bus.lastUpdated).toLocaleTimeString() : 'Unknown'}</p>
//                             </div>
//                         </Popup>
//                     </Marker>
//                 ))}
//             </MapContainer>
//         </div>
//     );
// };

// export default BusMap;




import React, { useState, useEffect } from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import L from 'leaflet';
import RoutePath from './RoutePath';
import { useWebSocket } from '../context/WebSocketContext';

// Fix for marker icons
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
    iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
});

// Custom bus icon
const busIcon = new L.DivIcon({
    className: 'bus-marker',
    html: '<i class="bi bi-bus-front"></i>',
    iconSize: [40, 40],
    iconAnchor: [20, 20]
});

// Custom stop icon
const stopIcon = new L.DivIcon({
    className: 'stop-marker',
    html: '<i class="bi bi-geo-alt-fill"></i>',
    iconSize: [20, 20],
    iconAnchor: [10, 10]
});

function MapUpdater({ center, zoom }) {
    const map = useMap();
    useEffect(() => {
        if (center) {
            map.setView(center, zoom);
        }
    }, [center, zoom, map]);
    return null;
}

const BusMap = ({ buses = [], stops = [], route = null, center = [37.7749, -122.4194], zoom = 13 }) => {
    const { buses: liveBuses } = useWebSocket();
    const [mapBuses, setMapBuses] = useState(buses);
    const [mapCenter, setMapCenter] = useState(center);
    const [mapZoom, setMapZoom] = useState(zoom);

    useEffect(() => {
        if (liveBuses.length > 0) {
            setMapBuses(prevBuses => {
                const updatedBuses = [...prevBuses].map(bus => {
                    const liveBus = liveBuses.find(lb => lb.id === bus.id);
                    return liveBus ? { ...bus, ...liveBus } : bus;
                });
                return updatedBuses;
            });
        }
    }, [liveBuses]);

    useEffect(() => {
        if (buses.length === 0 && liveBuses.length > 0) {
            setMapBuses(liveBuses);
        } else if (buses.length > 0) {
            setMapBuses(buses);
        }
    }, [buses, liveBuses]);

    useEffect(() => {
        if (route?.path?.length > 0) {
            const midPoint = Math.floor(route.path.length / 2);
            setMapCenter([route.path[midPoint].latitude, route.path[midPoint].longitude]);
            setMapZoom(14);
        } else if (stops.length > 0) {
            setMapCenter([stops[0].latitude, stops[0].longitude]);
            setMapZoom(15);
        }
    }, [route, stops]);

    return (
        <div className="map-container">
            <MapContainer center={mapCenter} zoom={mapZoom} style={{ height: '100%', width: '100%' }}>
                <MapUpdater center={mapCenter} zoom={mapZoom} />
                <TileLayer
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />

                {route && <RoutePath path={route.path} color={route.color} />}

                {stops.map(stop => (
                    <Marker
                        key={`stop-${stop.id}`}
                        position={[stop.latitude, stop.longitude]}
                        icon={stopIcon}
                    >
                        <Popup>
                            <div>
                                <h6>{stop.name}</h6>
                                <p>Stop #{stop.order}</p>
                                {route && <p>Route: {route.name}</p>}
                            </div>
                        </Popup>
                    </Marker>
                ))}

                {mapBuses.map(bus => (
                    <Marker
                        key={`bus-${bus.id}`}
                        position={[bus.latitude, bus.longitude]}
                        icon={busIcon}
                    >
                        <Popup>
                            <div>
                                <h6>Bus #{bus.busNumber}</h6>
                                <p>Route: {bus.route?.name || 'Unknown'}</p>
                                <p>Status: <span className={`badge status-${bus.status?.toLowerCase().replace(/\s/g, '-')}`}>
                                    {bus.status || 'Unknown'}
                                </span></p>
                                <p>Last Updated: {bus.lastUpdated ? new Date(bus.lastUpdated).toLocaleTimeString() : 'Unknown'}</p>
                            </div>
                        </Popup>
                    </Marker>
                ))}
            </MapContainer>
        </div>
    );
};

export default BusMap;