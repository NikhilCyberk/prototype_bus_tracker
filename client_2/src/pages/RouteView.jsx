// import React, { useState, useEffect } from 'react';
// import { useParams } from 'react-router-dom';
// import { Breadcrumb } from 'react-bootstrap';
// import { Link } from 'react-router-dom';
// import BusMap from '../components/BusMap';
// import RouteDetails from '../components/RouteDetails';
// import { getRouteById, getRouteStops, getBuses } from '../services/api';

// const RouteView = () => {
//     const { id } = useParams();
//     const [route, setRoute] = useState(null);
//     const [stops, setStops] = useState([]);
//     const [buses, setBuses] = useState([]);
//     const [loading, setLoading] = useState(true);
//     const [error, setError] = useState('');

//     useEffect(() => {
//         const fetchRouteData = async () => {
//             setLoading(true);
//             try {
//                 // Get route details
//                 const routeData = await getRouteById(id);
//                 setRoute(routeData);

//                 // Get stops for this route
//                 const stopsData = await getRouteStops(id);
//                 setStops(stopsData);

//                 // Get all buses and filter for this route
//                 const busesData = await getBuses();
//                 setBuses(busesData.filter(bus => bus.routeId === parseInt(id)));
//             } catch (error) {
//                 console.error('Error fetching route:', error);
//                 setError('Failed to load route data');
//             } finally {
//                 setLoading(false);
//             }
//         };

//         fetchRouteData();
//     }, [id]);

//     if (loading) {
//         return (
//             <div className="text-center py-5">
//                 <div className="spinner-border text-primary" role="status">
//                     <span className="visually-hidden">Loading...</span>
//                 </div>
//                 <p className="mt-3">Loading route information...</p>
//             </div>
//         );
//     }

//     if (error || !route) {
//         return (
//             <div className="text-center py-5">
//                 <div className="alert alert-danger" role="alert">
//                     {error || 'Route not found'}
//                 </div>
//                 <Link to="/" className="btn btn-primary mt-3">Back to Home</Link>
//             </div>
//         );
//     }

//     // Calculate map center based on route path or stops
//     const getMapCenter = () => {
//         if (route.path && route.path.length > 0) {
//             // Use the center of the route path
//             return [route.path[Math.floor(route.path.length / 2)].latitude,
//             route.path[Math.floor(route.path.length / 2)].longitude];
//         } else if (stops.length > 0) {
//             // Use the first stop
//             return [stops[0].latitude, stops[0].longitude];
//         }
//         // Default center
//         return [37.7749, -122.4194];
//     };

//     return (
//         <div className="route-view">
//             <Breadcrumb className="mb-4">
//                 <Breadcrumb.Item linkAs={Link} linkProps={{ to: "/" }}>Home</Breadcrumb.Item>
//                 <Breadcrumb.Item active>Route {route.routeNumber}</Breadcrumb.Item>
//             </Breadcrumb>

//             <h1 className="mb-4">
//                 <span
//                     className="route-color-indicator d-inline-block me-3"
//                     style={{
//                         width: '30px',
//                         height: '30px',
//                         borderRadius: '50%',
//                         backgroundColor: route.color,
//                         verticalAlign: 'middle'
//                     }}
//                 ></span>
//                 {route.name}
//             </h1>

//             <div className="mb-5">
//                 <BusMap
//                     buses={buses}
//                     stops={stops}
//                     route={route}
//                     center={getMapCenter()}
//                     zoom={14}
//                 />
//             </div>

//             <RouteDetails route={route} />
//         </div>
//     );
// };

// export default RouteView;




import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Breadcrumb, Row, Col, Card, Table, Badge, Spinner, Alert } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import BusMap from '../components/BusMap';
import LiveTracker from '../components/LiveTracker';
import { getRouteById, getRouteStops, getBuses } from '../services/api';
import { formatLastUpdated } from '../utils/helpers';

const RouteView = () => {
    const { id } = useParams();
    const [route, setRoute] = useState(null);
    const [stops, setStops] = useState([]);
    const [buses, setBuses] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchRouteData = async () => {
            setLoading(true);
            try {
                const [routeData, stopsData, busesData] = await Promise.all([
                    getRouteById(id),
                    getRouteStops(id),
                    getBuses()
                ]);

                setRoute(routeData);
                setStops(stopsData);
                setBuses(busesData.filter(bus => bus.routeId === parseInt(id)));
            } catch (error) {
                console.error('Error fetching route:', error);
                setError('Failed to load route data');
            } finally {
                setLoading(false);
            }
        };

        fetchRouteData();
    }, [id]);

    if (loading) {
        return (
            <div className="text-center py-5">
                <Spinner animation="border" variant="primary" />
                <p className="mt-3">Loading route information...</p>
            </div>
        );
    }

    if (error || !route) {
        return (
            <div className="text-center py-5">
                <Alert variant="danger">{error || 'Route not found'}</Alert>
                <Link to="/" className="btn btn-primary mt-3">Back to Home</Link>
            </div>
        );
    }

    return (
        <div className="route-view">
            <Breadcrumb className="mb-4">
                <Breadcrumb.Item linkAs={Link} linkProps={{ to: "/" }}>Home</Breadcrumb.Item>
                <Breadcrumb.Item active>Route {route.routeNumber}</Breadcrumb.Item>
            </Breadcrumb>

            <Row className="mb-4">
                <Col md={8}>
                    <h1 className="mb-4">
                        <span
                            className="route-color-indicator d-inline-block me-3"
                            style={{
                                width: '30px',
                                height: '30px',
                                borderRadius: '50%',
                                backgroundColor: route.color,
                                verticalAlign: 'middle'
                            }}
                        ></span>
                        {route.name}
                    </h1>

                    <Card className="mb-4">
                        <Card.Header className="bg-primary text-white">
                            <h5 className="mb-0">Route Map</h5>
                        </Card.Header>
                        <Card.Body className="p-0" style={{ height: '400px' }}>
                            <BusMap
                                buses={buses}
                                stops={stops}
                                route={route}
                            />
                        </Card.Body>
                    </Card>
                </Col>

                <Col md={4}>
                    <LiveTracker maxItems={5} />

                    <Card className="mt-4">
                        <Card.Header className="bg-secondary text-white">
                            <h5 className="mb-0">Quick Stats</h5>
                        </Card.Header>
                        <Card.Body>
                            <Table borderless size="sm">
                                <tbody>
                                    <tr>
                                        <td><strong>Route Number:</strong></td>
                                        <td>{route.routeNumber}</td>
                                    </tr>
                                    <tr>
                                        <td><strong>Stops:</strong></td>
                                        <td>{stops.length}</td>
                                    </tr>
                                    <tr>
                                        <td><strong>Active Buses:</strong></td>
                                        <td>
                                            {buses.filter(b => b.status === 'on-route').length} / {buses.length}
                                        </td>
                                    </tr>
                                    <tr>
                                        <td><strong>Color:</strong></td>
                                        <td>
                                            <span
                                                className="d-inline-block"
                                                style={{
                                                    width: '20px',
                                                    height: '20px',
                                                    backgroundColor: route.color,
                                                    verticalAlign: 'middle'
                                                }}
                                            ></span>
                                        </td>
                                    </tr>
                                </tbody>
                            </Table>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>

            <Row>
                <Col md={6}>
                    <Card className="mb-4">
                        <Card.Header className="bg-secondary text-white">
                            <h5 className="mb-0">Bus Stops</h5>
                        </Card.Header>
                        <Card.Body>
                            <Table striped hover responsive>
                                <thead>
                                    <tr>
                                        <th>#</th>
                                        <th>Stop Name</th>
                                        <th>Location</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {stops.map(stop => (
                                        <tr key={stop.id}>
                                            <td>{stop.order}</td>
                                            <td>
                                                <Link to={`/stops/${stop.id}`}>{stop.name}</Link>
                                            </td>
                                            <td>
                                                {stop.latitude.toFixed(4)}, {stop.longitude.toFixed(4)}
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </Table>
                        </Card.Body>
                    </Card>
                </Col>

                <Col md={6}>
                    <Card className="mb-4">
                        <Card.Header className="bg-secondary text-white">
                            <h5 className="mb-0">Active Buses</h5>
                        </Card.Header>
                        <Card.Body>
                            <Table striped hover responsive>
                                <thead>
                                    <tr>
                                        <th>Bus #</th>
                                        <th>Status</th>
                                        <th>Last Updated</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {buses.map(bus => (
                                        <tr key={bus.id}>
                                            <td>
                                                <Link to={`/buses/${bus.id}`}>{bus.busNumber}</Link>
                                            </td>
                                            <td>
                                                <Badge bg={
                                                    bus.status === 'on-route' ? 'success' :
                                                        bus.status === 'delayed' ? 'warning' :
                                                            bus.status === 'off-duty' ? 'danger' : 'secondary'
                                                }>
                                                    {bus.status || 'Unknown'}
                                                </Badge>
                                            </td>
                                            <td>
                                                {formatLastUpdated(bus.lastUpdated)}
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </Table>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default RouteView;