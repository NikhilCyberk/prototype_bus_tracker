import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Button, Form, Alert, ListGroup } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import BusMap from '../components/BusMap';
import { useAuth } from '../context/AuthContext';
import { useWebSocket } from '../context/WebSocketContext';
import { getBuses, getRoutes, updateBusStatus, updateBusLocation } from '../services/api';

const Dashboard = () => {
    const { user } = useAuth();
    const { updateBusStatus: wsUpdateStatus, updateBusLocation: wsUpdateLocation } = useWebSocket();
    const [buses, setBuses] = useState([]);
    const [routes, setRoutes] = useState([]);
    const [loading, setLoading] = useState(true);
    const [selectedBus, setSelectedBus] = useState(null);
    const [statusValue, setStatusValue] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    const isDriver = user?.role === 'driver';

    useEffect(() => {
        // const fetchData = async () => {
        //     try {
        //         const [busesData, routesData] = await Promise.all([
        //             getBuses(),
        //             getRoutes()
        //         ]);

        //         setBuses(busesData);
        //         setRoutes(routesData);

        //         // If user is a driver, find their bus
        //         if (isDriver) {
        //             const driverBus = busesData.find(bus => bus.driverId === user.id);
        //             setSelectedBus(driverBus || null);
        //             if (driverBus) {
        //                 setStatusValue(driverBus.status || 'off-duty');
        //             }
        //         }
        //     } catch (error) {
        //         console.error('Error fetching data:', error);
        //         setError('Failed to load data');
        //     } finally {
        //         setLoading(false);
        //     }
        // };

        // src/pages/Dashboard.jsx
        const fetchData = async () => {
            setLoading(true);
            setError('');
            try {
                const [busesData, routesData] = await Promise.all([
                    getBuses().catch(() => []), // Return empty array if error
                    getRoutes().catch(() => []) // Return empty array if error
                ]);

                setBuses(busesData || []);
                setRoutes(routesData || []);

                if (isDriver) {
                    const driverBus = busesData?.find(bus => bus.driverId === user.id);
                    setSelectedBus(driverBus || null);
                    if (driverBus) {
                        setStatusValue(driverBus.status || 'off-duty');
                    }
                }
            } catch (error) {
                console.error('Error fetching data:', error);
                setError('Failed to load data. Please try again later.');
                // Optionally set mock data here for development
                setBuses([]);
                setRoutes([]);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [user, isDriver]);

    const handleStatusChange = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        if (!selectedBus) {
            setError('No bus selected');
            return;
        }

        try {
            await updateBusStatus(selectedBus.id, statusValue);
            wsUpdateStatus(selectedBus.id, statusValue);
            setSuccess('Bus status updated successfully');

            // Update local state
            setBuses(prev =>
                prev.map(bus =>
                    bus.id === selectedBus.id ? { ...bus, status: statusValue } : bus
                )
            );
        } catch (error) {
            console.error('Error updating status:', error);
            setError('Failed to update bus status');
        }
    };

    const shareLocation = async () => {
        setError('');
        setSuccess('');

        if (!selectedBus) {
            setError('No bus selected');
            return;
        }

        if (!navigator.geolocation) {
            setError('Geolocation is not supported by your browser');
            return;
        }

        navigator.geolocation.getCurrentPosition(
            async (position) => {
                const { latitude, longitude } = position.coords;

                try {
                    await updateBusLocation(selectedBus.id, latitude, longitude);
                    wsUpdateLocation(selectedBus.id, latitude, longitude);
                    setSuccess('Location shared successfully');

                    // Update local state
                    setBuses(prev =>
                        prev.map(bus =>
                            bus.id === selectedBus.id ?
                                { ...bus, latitude, longitude, lastUpdated: new Date() } :
                                bus
                        )
                    );
                } catch (error) {
                    console.error('Error updating location:', error);
                    setError('Failed to update location');
                }
            },
            (err) => {
                console.error('Error getting location:', err);
                setError(`Failed to get location: ${err.message}`);
            },
            { enableHighAccuracy: true }
        );
    };

    if (loading) {
        return (
            <div className="text-center py-5">
                <div className="spinner-border text-primary" role="status">
                    <span className="visually-hidden">Loading...</span>
                </div>
                <p className="mt-3">Loading dashboard...</p>
            </div>
        );
    }

    return (
        <div className="dashboard">
            <h1 className="mb-4">Dashboard</h1>

            {isDriver ? (
                <div className="driver-dashboard">
                    <Row className="mb-4">
                        <Col>
                            <Card>
                                <Card.Header className="bg-primary text-white">
                                    <h4 className="mb-0">Driver Controls</h4>
                                </Card.Header>
                                <Card.Body>
                                    {error && <Alert variant="danger">{error}</Alert>}
                                    {success && <Alert variant="success">{success}</Alert>}

                                    {selectedBus ? (
                                        <>
                                            <h5>Bus #{selectedBus.busNumber}</h5>
                                            <p>Route: {routes.find(r => r.id === selectedBus.routeId)?.name || 'Unknown'}</p>

                                            <Form onSubmit={handleStatusChange} className="mb-4">
                                                <Form.Group className="mb-3">
                                                    <Form.Label>Bus Status</Form.Label>
                                                    <Form.Select
                                                        value={statusValue}
                                                        onChange={(e) => setStatusValue(e.target.value)}
                                                    >
                                                        <option value="on-route">On Route</option>
                                                        <option value="delayed">Delayed</option>
                                                        <option value="off-duty">Off Duty</option>
                                                    </Form.Select>
                                                </Form.Group>
                                                <Button type="submit" variant="primary">Update Status</Button>
                                            </Form>

                                            <Button
                                                variant="success"
                                                onClick={shareLocation}
                                                className="w-100"
                                            >
                                                Share Current Location
                                            </Button>
                                        </>
                                    ) : (
                                        <div className="text-center py-3">
                                            <p>No bus assigned to you.</p>
                                            <p>Please contact the administrator to assign a bus.</p>
                                        </div>
                                    )}
                                </Card.Body>
                            </Card>
                        </Col>
                    </Row>
                </div>
            ) : (
                <div className="passenger-dashboard">
                    <Row className="mb-4">
                        <Col md={8}>
                            <h2 className="mb-3">Live Bus Map</h2>
                            <BusMap buses={buses} />
                        </Col>
                        <Col md={4}>
                            <h2 className="mb-3">Quick Links</h2>
                            <Card>
                                <Card.Header className="bg-primary text-white">
                                    <h5 className="mb-0">Navigate</h5>
                                </Card.Header>
                                <ListGroup variant="flush">
                                    <ListGroup.Item action as={Link} to="/">
                                        <i className="bi bi-house-door me-2"></i> Home
                                    </ListGroup.Item>
                                    <ListGroup.Item action as="button" onClick={() => window.location.reload()}>
                                        <i className="bi bi-arrow-clockwise me-2"></i> Refresh Data
                                    </ListGroup.Item>
                                </ListGroup>
                            </Card>
                        </Col>
                    </Row>

                    <h2 className="mb-3">Available Routes</h2>
                    <Row xs={1} md={2} lg={3} className="g-4 mb-4">
                        {routes.map(route => (
                            <Col key={route.id}>
                                <Card className="h-100 dashboard-card">
                                    <Card.Header style={{ backgroundColor: route.color, color: '#fff' }}>
                                        Route {route.routeNumber}
                                    </Card.Header>
                                    <Card.Body>
                                        <Card.Title>{route.name}</Card.Title>
                                        <Card.Text>
                                            {route.stops?.length || 0} stops along the route
                                        </Card.Text>
                                    </Card.Body>
                                    <Card.Footer>
                                        <Button
                                            as={Link}
                                            to={`/routes/${route.id}`}
                                            variant="outline-primary"
                                            size="sm"
                                            className="w-100"
                                        >
                                            View Route Details
                                        </Button>
                                    </Card.Footer>
                                </Card>
                            </Col>
                        ))}
                    </Row>
                </div>
            )}
        </div>
    );
};

export default Dashboard;