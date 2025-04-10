import React, { useState, useEffect } from 'react';
import { Card, ListGroup, Badge } from 'react-bootstrap';
import { getStopArrivals } from '../services/api';
import { formatTimeRemaining } from '../utils/helpers';

const StopInfo = ({ stop }) => {
    const [arrivals, setArrivals] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchArrivals = async () => {
            setLoading(true);
            try {
                const data = await getStopArrivals(stop.id);
                setArrivals(data);
            } catch (error) {
                console.error('Error fetching arrivals:', error);
            } finally {
                setLoading(false);
            }
        };

        if (stop?.id) {
            fetchArrivals();

            // Refresh arrivals every 30 seconds
            const interval = setInterval(fetchArrivals, 30000);
            return () => clearInterval(interval);
        }
    }, [stop]);

    const calculateTimeRemaining = (estimatedTime) => {
        const now = Math.floor(Date.now() / 1000);
        const diff = estimatedTime - now;

        if (diff <= 0) return 'Arriving now';
        if (diff < 60) return `${diff} seconds`;

        const minutes = Math.floor(diff / 60);
        return `${minutes} min${minutes !== 1 ? 's' : ''}`;
    };

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
        <Card className="stop-arrivals">
            <Card.Header className="bg-primary text-white">
                <h5 className="mb-0">Upcoming Arrivals</h5>
            </Card.Header>
            <ListGroup variant="flush">
                {arrivals.length > 0 ? (
                    arrivals.map((arrival, index) => (
                        <ListGroup.Item
                            key={index}
                            className="d-flex justify-content-between align-items-center"
                        >
                            <div>
                                <h6 className="mb-0">Bus #{arrival.busNumber}</h6>
                                <small className="text-muted">Route: {arrival.routeNumber}</small>
                            </div>
                            {/* <Badge bg="primary" pill>
                                {calculateTimeRemaining(arrival.estimatedTime)}
                            </Badge> */}

                            <Badge bg="primary" pill>
                                {formatTimeRemaining(arrival.estimatedTime - Math.floor(Date.now() / 1000))}
                            </Badge>
                        </ListGroup.Item>
                    ))
                ) : (
                    <ListGroup.Item className="text-center py-3">
                        <p className="mb-0">No upcoming arrivals</p>
                    </ListGroup.Item>
                )}
            </ListGroup>
        </Card>
    );
};

export default StopInfo;