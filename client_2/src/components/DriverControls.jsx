import React, { useState, useEffect } from 'react';
import { Button, Form, Alert, Card } from 'react-bootstrap';
import { useWebSocket } from '../context/WebSocketContext';
import { getBuses } from '../services/api';

const DriverControls = ({ user }) => {
    const { updateBusLocation, updateBusStatus } = useWebSocket();
    const [bus, setBus] = useState(null);
    const [status, setStatus] = useState('off-duty');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    useEffect(() => {
        const fetchDriverBus = async () => {
            try {
                const buses = await getBuses();
                const driverBus = buses.find(b => b.driverId === user.id);
                setBus(driverBus);
                if (driverBus) {
                    setStatus(driverBus.status || 'off-duty');
                }
            } catch (err) {
                setError('Failed to load bus information');
            } finally {
                setLoading(false);
            }
        };

        fetchDriverBus();
    }, [user]);

    const handleStatusUpdate = async (e) => {
        e.preventDefault();
        if (!bus) return;

        try {
            await updateBusStatus(bus.id, status);
            setSuccess('Status updated successfully');
            setError('');
        } catch (err) {
            setError('Failed to update status');
            setSuccess('');
        }
    };

    const handleLocationUpdate = async () => {
        if (!bus) return;

        if (!navigator.geolocation) {
            setError('Geolocation is not supported by your browser');
            return;
        }

        navigator.geolocation.getCurrentPosition(
            async (position) => {
                try {
                    await updateBusLocation(
                        bus.id,
                        position.coords.latitude,
                        position.coords.longitude
                    );
                    setSuccess('Location updated successfully');
                    setError('');
                } catch (err) {
                    setError('Failed to update location');
                    setSuccess('');
                }
            },
            (err) => {
                setError(`Geolocation error: ${err.message}`);
            },
            { enableHighAccuracy: true, timeout: 10000 }
        );
    };

    if (loading) {
        return <div className="text-center py-3">Loading bus information...</div>;
    }

    if (!bus) {
        return (
            <Alert variant="warning">
                No bus assigned to you. Please contact the administrator.
            </Alert>
        );
    }

    return (
        <Card>
            <Card.Header className="bg-primary text-white">
                <h5>Bus #{bus.busNumber} Controls</h5>
            </Card.Header>
            <Card.Body>
                {error && <Alert variant="danger">{error}</Alert>}
                {success && <Alert variant="success">{success}</Alert>}

                <Form onSubmit={handleStatusUpdate}>
                    <Form.Group className="mb-3">
                        <Form.Label>Current Status</Form.Label>
                        <Form.Select
                            value={status}
                            onChange={(e) => setStatus(e.target.value)}
                        >
                            <option value="on-route">On Route</option>
                            <option value="delayed">Delayed</option>
                            <option value="off-duty">Off Duty</option>
                        </Form.Select>
                    </Form.Group>
                    <Button type="submit" variant="primary" className="w-100 mb-3">
                        Update Status
                    </Button>
                </Form>

                <Button
                    variant="success"
                    onClick={handleLocationUpdate}
                    className="w-100"
                >
                    Share Current Location
                </Button>
            </Card.Body>
        </Card>
    );
};

export default DriverControls;