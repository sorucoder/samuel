import React, { useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import Alert from 'react-bootstrap/Alert';

import SAMUELAlertLink from './SAMUELAlertLink.js';

function SAMUELAlert({ variant, show, redirect = null, children }) {
    const navigate = useNavigate();
    useEffect(() => {
        if (redirect) {
            setTimeout(() => navigate(redirect.path), 5000);
        }
    }, [redirect]);

    if (redirect) {
        return (
            <Alert variant={variant} show={show}>
                {children}
                <hr />
                You will be redirected to the <Alert.Link as={Link} to={redirect.path}>{redirect.name}</Alert.Link> in 5 seconds...
            </Alert>
        );
    } else {
        return (
            <Alert variant={variant} show={show}>{children}</Alert>
        );
    }
}

SAMUELAlert.Link = SAMUELAlertLink;

export default SAMUELAlert;