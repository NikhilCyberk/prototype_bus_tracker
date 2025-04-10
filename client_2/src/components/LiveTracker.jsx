// import React, { useState, useEffect } from 'react';
// import { Card, ListGroup, Badge } from 'react-bootstrap';
// import { Link } from 'react-router-dom';
// import { getBuses } from '../services/api';
// import { useWebSocket } from '../context/WebSocketContext';

// const LiveTracker = () => {
//     const [buses, setBuses] = useState([]);
//     const [loading, setLoading] = useState(true);
//     const { buses: liveBuses, connected } = useWebSocket();

//     useEffect(() => {
//         const fetchBuses = async () => {
//             try {
//                 const data = await getBuses();
//                 setBuses(data);
//             } catch (error) {
//                 console.error('Error fetching buses:', error);
//             } finally {
//                 setLoading(false);
//             }
//         };

//         fetchBuses();
//     }, []);

//     // Update buses with real-time data
//     useEffect(() => {
//         if (liveBuses.length > 0) {
//             setBuses(prevBuses => {
//                 const updatedBuses = [...prevBuses];

//                 liveBuses.forEach(liveBus => {
//                     const index = updatedBuses.findIndex(b => b.id === liveBus.id);
//                     if (index !== -1) {
//                         updatedBuses[index] = { ...updatedBuses[index], ...liveBus };
//                     }
//                 });

//                 return updatedBuses;
//             });
//         }
//     }, [liveBuses]);

//     const getStatusBadgeClass = (status) => {
//         switch (status?.toLowerCase()) {
//             case 'on-route':
//                 return 'success';
//             case 'delayed':
//                 return 'warning';
//             case 'off-duty':
//                 return 'danger';
//             default:
//                 return 'secondary';
//         }
//     };

//     if (loading) {
//         return (
//             <Card>
//                 <Card.Header className="bg-primary text-white">
//                     <h5 className="mb-0">Live Bus Status</h5>
//                 </Card.Header>
//                 <Card.Body className="text-center py-4">
//                     <div className="spinner-border text-primary" role="status">
//                         <span className="visually-hidden">Loading...</span>
//                     </div>
//                     <p className="mt-3">Loading bus information...</p>
//                 </Card.Body>
//             </Card>
//         );
//     }

//     return (
//         <Card>
//             <Card.Header className="bg-primary text-white d-flex justify-content-between align-items-center">
//                 <h5 className="mb-0">Live Bus Status</h5>
//                 <Badge bg={connected ? 'success' : 'danger'}>
//                     {connected ? 'Connected' : 'Disconnected'}
//                 </Badge>
//             </Card.Header>
//             <ListGroup variant="flush">
//                 {buses.length > 0 ? (
//                     buses.map(bus => (
//                         <ListGroup.Item key={bus.id} className="d-flex justify-content-between align-items-center">
//                             <div>
//                                 <h6><Link to={`/buses/${bus.id}`}>Bus #{bus.busNumber}</Link></h6>
//                                 <small>Route: {bus.route?.name || 'Unknown'}</small>
//                             </div>
//                             <div className="d-flex align-items-center">
//                                 <Badge bg={getStatusBadgeClass(bus.status)}>
//                                     {bus.status || 'Unknown'}
//                                 </Badge>
//                                 {bus.lastUpdated && (
//                                     <small className="ms-2 text-muted">
//                                         {new Date(bus.lastUpdated).toLocaleTimeString()}
//                                     </small>
//                                 )}
//                             </div>
//                         </ListGroup.Item>
//                     ))
//                 ) : (
//                     <ListGroup.Item className="text-center py-3">
//                         <p className="mb-0">No buses currently active</p>
//                     </ListGroup.Item>
//                 )}
//             </ListGroup>
//         </Card>
//     );
// };

// export default LiveTracker;



import React, { useState, useEffect } from 'react';
import { Card, ListGroup, Badge, Spinner } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { getBuses } from '../services/api';
import { useWebSocket } from '../context/WebSocketContext';
import { formatLastUpdated } from '../utils/helpers';

const LiveTracker = ({ maxItems = 10 }) => {
    const [buses, setBuses] = useState([]);
    const [loading, setLoading] = useState(true);
    const { buses: liveBuses, connected } = useWebSocket();

    useEffect(() => {
        const fetchBuses = async () => {
            try {
                const data = await getBuses();
                setBuses(data.slice(0, maxItems));
            } catch (error) {
                console.error('Error fetching buses:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchBuses();
    }, [maxItems]);

    useEffect(() => {
        if (liveBuses.length > 0) {
            setBuses(prevBuses => {
                const updatedBuses = [...prevBuses];

                liveBuses.forEach(liveBus => {
                    const index = updatedBuses.findIndex(b => b.id === liveBus.id);
                    if (index !== -1) {
                        updatedBuses[index] = { ...updatedBuses[index], ...liveBus };
                    }
                });

                return updatedBuses.slice(0, maxItems);
            });
        }
    }, [liveBuses, maxItems]);

    const getStatusBadgeClass = (status) => {
        if (!status) return 'secondary';
        const statusLower = status.toLowerCase();
        if (statusLower.includes('on-route')) return 'success';
        if (statusLower.includes('delayed')) return 'warning';
        if (statusLower.includes('off-duty')) return 'danger';
        return 'secondary';
    };

    if (loading) {
        return (
            <Card>
                <Card.Header className="bg-primary text-white">
                    <h5 className="mb-0">Live Bus Status</h5>
                </Card.Header>
                <Card.Body className="text-center py-4">
                    <Spinner animation="border" variant="primary" />
                    <p className="mt-3">Loading bus information...</p>
                </Card.Body>
            </Card>
        );
    }

    return (
        <Card>
            <Card.Header className="bg-primary text-white d-flex justify-content-between align-items-center">
                <h5 className="mb-0">Live Bus Status</h5>
                <Badge bg={connected ? 'success' : 'danger'}>
                    {connected ? 'Connected' : 'Disconnected'}
                </Badge>
            </Card.Header>
            <ListGroup variant="flush">
                {buses.length > 0 ? (
                    buses.map(bus => (
                        <ListGroup.Item key={bus.id} className="d-flex justify-content-between align-items-center">
                            <div>
                                <h6 className="mb-0">
                                    <Link to={`/buses/${bus.id}`} className="text-decoration-none">
                                        Bus #{bus.busNumber}
                                    </Link>
                                </h6>
                                <small className="text-muted">
                                    {bus.route?.name || 'Unknown Route'}
                                </small>
                            </div>
                            <div className="d-flex flex-column align-items-end">
                                <Badge bg={getStatusBadgeClass(bus.status)} className="mb-1">
                                    {bus.status || 'Unknown'}
                                </Badge>
                                <small className="text-muted">
                                    {formatLastUpdated(bus.lastUpdated)}
                                </small>
                            </div>
                        </ListGroup.Item>
                    ))
                ) : (
                    <ListGroup.Item className="text-center py-3">
                        <p className="mb-0">No buses currently active</p>
                    </ListGroup.Item>
                )}
                {buses.length > 0 && maxItems && (
                    <ListGroup.Item className="text-center">
                        <Link to="/dashboard" className="text-primary">
                            View all buses
                        </Link>
                    </ListGroup.Item>
                )}
            </ListGroup>
        </Card>
    );
};

export default LiveTracker;