import React, { useState, useEffect } from 'react';
import { Card, Button, Row, Col } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import BusMap from '../components/BusMap';
import LiveTracker from '../components/LiveTracker';
import { getRoutes, getBuses } from '../services/api';

const Home = () => {
    const [routes, setRoutes] = useState([]);
    const [buses, setBuses] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [routesData, busesData] = await Promise.all([
                    getRoutes(),
                    getBuses()
                ]);

                setRoutes(routesData);
                setBuses(busesData);
            } catch (error) {
                console.error('Error fetching data:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    return (
        <div className="home-page">
            <header className="text-center mb-5">
                <h1 className="display-4">Welcome to CityBus Tracker</h1>
                <p className="lead">Track real-time bus locations, schedules, and routes</p>
            </header>

            <div className="mb-5">
                <h2 className="mb-4">Live Bus Map</h2>
                <BusMap buses={buses} />
            </div>

            <Row className="mb-5">
                <Col md={7}>
                    <h2 className="mb-4">Available Routes</h2>
                    {loading ? (
                        <div className="text-center py-4">
                            <div className="spinner-border text-primary" role="status">
                                <span className="visually-hidden">Loading...</span>
                            </div>
                        </div>
                    ) : (
                        <Row xs={1} md={2} className="g-4">
                            {routes.map(route => (
                                <Col key={route.id}>
                                    <Card className="h-100 route-card">
                                        <Card.Header style={{ backgroundColor: route.color, color: '#fff' }}>
                                            Route {route.routeNumber}
                                        </Card.Header>
                                        <Card.Body>
                                            <Card.Title>{route.name}</Card.Title>
                                            <Card.Text>
                                                {route.stops?.length || 0} stops
                                            </Card.Text>
                                            <Button
                                                as={Link}
                                                to={`/routes/${route.id}`}
                                                variant="outline-primary"
                                                size="sm"
                                            >
                                                View Details
                                            </Button>
                                        </Card.Body>
                                    </Card>
                                </Col>
                            ))}
                        </Row>
                    )}
                </Col>
                <Col md={5}>
                    <h2 className="mb-4">Live Status</h2>
                    <LiveTracker />
                </Col>
            </Row>

            <div className="bg-light p-4 rounded-3 mb-5">
                <h2 className="mb-3">How to Use</h2>
                <Row xs={1} md={3} className="g-4">
                    <Col>
                        <div className="text-center">
                            <div className="display-5 mb-3">
                                <i className="bi bi-map text-primary"></i>
                            </div>
                            <h4>View Routes</h4>
                            <p>Browse available bus routes and check their schedules.</p>
                        </div>
                    </Col>
                    <Col>
                        <div className="text-center">
                            <div className="display-5 mb-3">
                                <i className="bi bi-geo-alt text-primary"></i>
                            </div>
                            <h4>Find Stops</h4>
                            <p>Locate nearby bus stops and check upcoming arrivals.</p>
                        </div>
                    </Col>
                    <Col>
                        <div className="text-center">
                            <div className="display-5 mb-3">
                                <i className="bi bi-clock-history text-primary"></i>
                            </div>
                            <h4>Real-Time Updates</h4>
                            <p>Get live updates on bus locations and estimated arrival times.</p>
                        </div>
                    </Col>
                </Row>
            </div>
        </div>
    );
};

export default Home;