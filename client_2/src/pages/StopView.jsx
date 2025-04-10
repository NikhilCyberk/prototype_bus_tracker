import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Breadcrumb, Row, Col, Card } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import BusMap from '../components/BusMap';
import StopInfo from '../components/StopInfo';
import { getStopById, getRouteById } from '../services/api';

const StopView = () => {
    const { id } = useParams();
    const [stop, setStop] = useState(null);
    const [route, setRoute] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchStopData = async () => {
            setLoading(true);
            try {
                // Get stop details
                const stopData = await getStopById(id);
                setStop(stopData);

                // Get route for this stop
                if (stopData.routeId) {
                    const routeData = await getRouteById(stopData.routeId);
                    setRoute(routeData);
                }
            } catch (error) {
                console.error('Error fetching stop:', error);
                setError('Failed to load stop data');
            } finally {
                setLoading(false);
            }
        };

        fetchStopData();
    }, [id]);

    if (loading) {
        return (
            <div className="text-center py-5">
                <div className="spinner-border text-primary" role="status">
                    <span className="visually-hidden">Loading...</span>
                </div>
                <p className="mt-3">Loading stop information...</p>
            </div>
        );
    }

    if (error || !stop) {
        return (
            <div className="text-center py-5">
                <div className="alert alert-danger" role="alert">
                    {error || 'Stop not found'}
                </div>
                <Link to="/" className="btn btn-primary mt-3">Back to Home</Link>
            </div>
        );
    }

    return (
        <div className="stop-view">
            <Breadcrumb className="mb-4">
                <Breadcrumb.Item linkAs={Link} linkProps={{ to: "/" }}>Home</Breadcrumb.Item>
                {route && (
                    <Breadcrumb.Item linkAs={Link} linkProps={{ to: `/routes/${route.id}` }}>
                        Route {route.routeNumber}
                    </Breadcrumb.Item>
                )}
                <Breadcrumb.Item active>Stop: {stop.name}</Breadcrumb.Item>
            </Breadcrumb>

            <h1 className="mb-4">{stop.name}</h1>

            <Row className="mb-4">
                <Col md={8}>
                    <Card className="mb-4">
                        <Card.Header className="bg-primary text-white">
                            <h4 className="mb-0">Stop Details</h4>
                        </Card.Header>
                        <Card.Body>
                            <Row>
                                <Col md={6}>
                                    <p><strong>Route:</strong> {route ? route.name : 'Unknown'}</p>
                                    <p><strong>Stop Order:</strong> {stop.order}</p>
                                </Col>
                                <Col md={6}>
                                    <p><strong>Location:</strong></p>
                                    <p>Latitude: {stop.latitude.toFixed(6)}</p>
                                    <p>Longitude: {stop.longitude.toFixed(6)}</p>
                                </Col>
                            </Row>
                        </Card.Body>
                    </Card>

                    <div className="map-container mb-4">
                        <BusMap
                            stops={[stop]}
                            center={[stop.latitude, stop.longitude]}
                            zoom={16}
                        />
                    </div>
                </Col>

                <Col md={4}>
                    <StopInfo stop={stop} />

                    {route && (
                        <Card className="mt-4">
                            <Card.Header className="bg-secondary text-white">
                                <h5 className="mb-0">Route Information</h5>
                            </Card.Header>
                            <Card.Body>
                                <h6>{route.name}</h6>
                                <p>Route Number: {route.routeNumber}</p>
                                <p>Color:
                                    <span
                                        className="d-inline-block ms-2"
                                        style={{
                                            width: '20px',
                                            height: '20px',
                                            backgroundColor: route.color,
                                            verticalAlign: 'middle'
                                        }}
                                    ></span>
                                </p>
                                <Button
                                    as={Link}
                                    to={`/routes/${route.id}`}
                                    variant="outline-primary"
                                    size="sm"
                                    className="w-100"
                                >
                                    View Full Route
                                </Button>
                            </Card.Body>
                        </Card>
                    )}
                </Col>
            </Row>
        </div>
    );
};

export default StopView;