import React from 'react';
import { Polyline } from 'react-leaflet';
import L from 'leaflet';

const RoutePath = ({ path, color = '#3388ff', dashArray = null }) => {
    if (!path || path.length < 2) return null;

    // Convert path to LatLng array
    const positions = path.map(point => [point.latitude, point.longitude]);

    // Create a custom dashed pattern if needed
    const pathOptions = {
        color,
        weight: 5,
        opacity: 0.7,
        lineCap: 'round',
        lineJoin: 'round'
    };

    if (dashArray) {
        pathOptions.dashArray = dashArray;
    }

    return <Polyline positions={positions} pathOptions={pathOptions} />;
};

export default RoutePath;