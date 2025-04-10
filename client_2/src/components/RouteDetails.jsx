import React, { useState, useEffect } from 'react';
import { Card, Table, Badge } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { getRouteStops, getBuses } from '../services/api';

const RouteDetails = ({ route }) => {
    const [stops, setStops] = useState([]);
    const [buses, setBuses] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            setLoading(true);
            try {
                // Get stops for this route
                const stopsData = await getRouteStops(route.id);
                setStops(stopsData);

                // Get all buses and filter for this route
                const busesData = await getBuses();
                setBuses(busesData.filter(bus => bus.routeId === route.id));
            } catch (error) {
                console.error('Error fetching route details:', error);
            } finally {
                setLoading(false);
            }
        };

        if (route?.id) {
            fetchData();
        }
    }, [route]);

    if (loading) {
        return (
            <div className="text-center py-4">
                <div className="spinner-border text-primary" role="status">
                    <span className="visually-hidden">Loading...</span>
                </div>
            </div>
        );
    }

    return (
        <div className="route-details">
            <Card className="mb-4">
                <Card.Header className="bg-primary text-white">
                    <h5 className="mb-0">Route Information</h5>
                </Card.Header>
                <Card.Body>
                    <div className="d-flex align-items-center mb-3">
                        <div
                            className="route-color-indicator me-3"
                            style={{
                                width: '30px',
                                height: '30px',
                                borderRadius: '50%',
                                backgroundColor: route.color
                            }}
                        ></div>
                        <div>
                            <h4 className="mb-0">{route.name}</h4>
                            <p className="text-muted mb-0">Route Number: {route.routeNumber}</p>
                        </div>
                    </div>
                </Card.Body>
            </Card>

            <div className="row">
                <div className="col-md-6 mb-4">
                    <Card>
                        <Card.Header className="bg-secondary text-white">
                            <h5 className="mb-0">Bus Stops</h5>
                        </Card.Header>
                        <Card.Body>
                            {stops.length > 0 ? (
                                <Table striped hover responsive>
                                    <thead>
                                        <tr>
                                            <th>Order</th>
                                            <th>Name</th>
                                            <th>Details</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {stops.map(stop => (
                                            <tr key={stop.id}>
                                                <td>{stop.order}</td>
                                                <td>{stop.name}</td>
                                                <td>
                                                    <Link to={`/stops/${stop.id}`}>View Stop</Link>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </Table>
                            ) : (
                                <p className="text-center mb-0">No stops found for this route</p>
                            )}
                        </Card.Body>
                    </Card>
                </div>

                <div className="col-md-6 mb-4">
                    <Card>
                        <Card.Header className="bg-secondary text-white">
                            <h5 className="mb-0">Buses on this Route</h5>
                        </Card.Header>
                        <Card.Body>
                            {buses.length > 0 ? (
                                <Table striped hover responsive>
                                    <thead>
                                        <tr>
                                            <th>Bus Number</th>
                                            <th>Status</th>
                                            <th>Last Updated</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {buses.map(bus => (
                                            <tr key={bus.id}>
                                                <td>{bus.busNumber}</td>
                                                <td>
                                                    <Badge
                                                        bg={bus.status === 'on-route' ? 'success' :
                                                            bus.status === 'delayed' ? 'warning' :
                                                                bus.status === 'off-duty' ? 'danger' : 'secondary'}
                                                    >
                                                        {bus.status || 'Unknown'}
                                                    </Badge>
                                                </td>
                                                <td>
                                                    {bus.lastUpdated ?
                                                        new Date(bus.lastUpdated).toLocaleTimeString() :
                                                        'N/A'}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </Table>
                            ) : (
                                <p className="text-center mb-0">No buses currently on this route</p>
                            )}
                        </Card.Body>
                    </Card>
                </div>
            </div>
        </div>
    );
};

export default RouteDetails;