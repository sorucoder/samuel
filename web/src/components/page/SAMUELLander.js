import React from 'react';

import Card from 'react-bootstrap/Card';
import CardGroup from 'react-bootstrap/CardGroup';
import Container from 'react-bootstrap/Container';

import mountainsImage from '../../images/mountains.png';
import logoImage from '../../images/logo.png';

export default function SAMUELLander({title, subtitle, children}) {
    return (
        <Container fluid id="main" style={{background: 'linear-gradient(90deg, rgba(2, 0, 36, 1) 0%, rgba(49, 113, 184, 1) 47%, rgba(0, 212, 255, 1) 87%)'}} className="min-vh-100 d-flex align-items-center justify-content-center">
            <CardGroup className="w-75 min-vh-50">
                <Card>
                    <Card.Img variant="bottom" src={mountainsImage} alt="Mountains" className="w-100 h-100 object-fit-cover" />
                    <Card.ImgOverlay className="text-white">
                        <img src={logoImage} id="logo" className="d-none d-md-inline" width="64px" height="64px" alt="South Hills Initials logo" />
                        <h1 className="position-absolute start-50 top-50 translate-middle text-center">{ title }</h1>
                    </Card.ImgOverlay>
                </Card>
                <Card>
                    <Card.Title as="h2" className="my-3 text-center">{ subtitle }</Card.Title>
                    <Card.Body className="d-flex flex-column">
                        { children }
                    </Card.Body>
                </Card>
            </CardGroup>
        </Container>
    );
}