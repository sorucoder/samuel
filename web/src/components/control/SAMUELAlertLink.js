import React from 'react';
import { Link as RouterLink } from 'react-router-dom';

import AlertLink from 'react-bootstrap/AlertLink';

export default function SAMUELAlertLink({to, children}) {
    return (
        <AlertLink as={RouterLink} to={to}>{children}</AlertLink>
    );
}