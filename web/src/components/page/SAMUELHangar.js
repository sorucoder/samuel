import React, { useEffect, useState } from 'react';
import { NavLink as RouterNavLink, redirect } from 'react-router-dom';
import { Duration } from 'luxon';

import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import Modal from 'react-bootstrap/Modal';
import Nav from 'react-bootstrap/Nav';
import Offcanvas from 'react-bootstrap/Offcanvas';

import {
    PeopleFill as PeopleIcon,
    Bank2 as CampusIcon,
    BookFill as ProgramIcon,
    BuildingsFill as CompanyIcon,
    ListCheck as AuditIcon,
    CalendarWeekFill as TimecardIcon,
    FileEarmarkTextFill as ReportIcon,
    BellFill as NotificationIcon,
    DoorOpenFill as LogoutIcon,
    List as NavigatorMenuIcon
} from 'react-bootstrap-icons';

const NAVIGATOR_ICON_RENDER_PROPS = {
    size: '1.25rem',
    className: 'me-2'
}

function NavigatorLink({ to, last = false, children }) {
    if (last) {
        return (
            <Nav.Item className="mt-auto">
                <Nav.Link as={RouterNavLink} to={to}>
                    {children}
                </Nav.Link>
            </Nav.Item>
        );
    } else {
        return (
            <Nav.Item>
                <Nav.Link as={RouterNavLink} to={to}>
                    {children}
                </Nav.Link>
            </Nav.Item>
        );
    }
}

const NAVIGATOR_LINKS = [
    {
        roles: ['administrator'],
        element: <NavigatorLink key="instructors" to="/instructor/view">
            <PeopleIcon {...NAVIGATOR_ICON_RENDER_PROPS} />
            Instructors
        </NavigatorLink>
    },
    {
        roles: ['administrator', 'instructor'],
        element: <NavigatorLink key="supervisors" to="/supervisor/view">
            <PeopleIcon {...NAVIGATOR_ICON_RENDER_PROPS} />
            Supervisors
        </NavigatorLink>
    },
    {
        roles: ['administrator', 'instructor', 'supervisor'],
        element: <NavigatorLink key="students" to="/student/view">
            <PeopleIcon {...NAVIGATOR_ICON_RENDER_PROPS} />
            Students
        </NavigatorLink>
    },
    {
        roles: ['administrator'],
        element: <NavigatorLink key="campuses" to="/campus/view">
            <CampusIcon {...NAVIGATOR_ICON_RENDER_PROPS} />
            Campuses
        </NavigatorLink>
    },
    {
        roles: ['administrator'],
        element: <NavigatorLink key="programs" to="/program/view">
            <ProgramIcon {...NAVIGATOR_ICON_RENDER_PROPS} />
            Programs
        </NavigatorLink>
    },
    {
        roles: ['administrator', 'instructor'],
        element: <NavigatorLink key="companies" to="/company/view">
            <CompanyIcon {...NAVIGATOR_ICON_RENDER_PROPS} />
            Companies
        </NavigatorLink>
    },
    {
        roles: ['administrator', 'instructor', 'supervisor', 'student'],
        element: <NavigatorLink key="timecards" to="/timecard/view">
            <TimecardIcon {...NAVIGATOR_ICON_RENDER_PROPS} />
            Timecards
        </NavigatorLink>
    },
    {
        roles: ['administrator', 'instructor', 'supervisor', 'student'],
        element: <NavigatorLink key="reports" to="/report/view">
            <ReportIcon {...NAVIGATOR_ICON_RENDER_PROPS} />
            Reports
        </NavigatorLink>
    },
    {
        roles: ['administrator'],
        element: <NavigatorLink key="audit" to="/audit/view">
            <AuditIcon {...NAVIGATOR_ICON_RENDER_PROPS} />
            Audit
        </NavigatorLink>
    },
    {
        roles: ['administrator', 'instructor', 'supervisor', 'student'],
        element: <NavigatorLink key="notifications" to="/notification/view">
            <NotificationIcon {...NAVIGATOR_ICON_RENDER_PROPS} />
            Notifications
        </NavigatorLink>
    },
    {
        roles: ['administrator', 'instructor', 'supervisor', 'student'],
        element: <NavigatorLink key="logout" to="/logout" last>
            <LogoutIcon {...NAVIGATOR_ICON_RENDER_PROPS} />
            Logout
        </NavigatorLink>
    }
];

function Navigator({show, onHide, placement, user}) {
    return (
        <Offcanvas show={show} onHide={onHide} placement={placement} className="text-bg-primary">
            <Offcanvas.Header closeButton >
                <Nav.Link as={RouterNavLink} to="/dashboard">
                    <Offcanvas.Title as="h2" className="d-flex justify-content-center gap-2">
                        Navigation
                    </Offcanvas.Title>
                </Nav.Link>
            </Offcanvas.Header>
            <Offcanvas.Body className="text-bg-dark">
                <Nav className="flex-column h-100">
                    {NAVIGATOR_LINKS.map((navigatorLink) => {
                        if (navigatorLink.roles.includes(user.role.id)) {
                            return navigatorLink.element;
                        }
                    })}
                </Nav>
            </Offcanvas.Body>
        </Offcanvas>
    );
}

const SESSION_EXPIRATION_WARNING = Duration.fromObject({minutes: 2}).shiftTo('milliseconds');

function SessionExpirationModal({session}) {
    const [show, setShow] = useState(false);
    const [untilExpiration, setUntilExpiration] = useState(
        session.expiresOn
            .diffNow(['minutes', 'seconds'])
            .mapUnits((value) => Math.floor(value))
    );

    useEffect(() => {
        const untilExpiration = session.expiresOn.diffNow();
        
        let countdownInterval;
        const warnTimeout = setTimeout(() => {
            setShow(true);
            countdownInterval = setInterval(() => {
                setUntilExpiration(
                    session.expiresOn
                        .diffNow(['minutes', 'seconds'])
                        .mapUnits((value) => Math.floor(value))
                );
            }, 1000);
        }, Math.max(0, untilExpiration.milliseconds - SESSION_EXPIRATION_WARNING.milliseconds));
        const expireTimeout = setTimeout(() => redirect('/login?expired=true'), untilExpiration.milliseconds);

        return () => {
            clearTimeout(warnTimeout);
            if (countdownInterval) {
                clearInterval(countdownInterval);
            }
            clearTimeout(expireTimeout);
        }
    }, [session.expiresOn]);

    return (
        <Modal show={show}>
            <Modal.Header closeButton>Your Session Is About To Expire</Modal.Header>
            <Modal.Body>
                Your session is set to expire in {untilExpiration.toHuman()}.
                Dismiss this modal to continue your session.
            </Modal.Body>
        </Modal>
    );
}

function Header({user}) {
    const [showNavigator, setShowNavigator] = useState(false);

    return (
        <>
            <Navigator show={showNavigator} onHide={() => setShowNavigator(false)} placement="end" user={user} />
            <Container fluid className="d-flex flex-row justify-content-between align-items-center text-bg-primary p-3">
                <hgroup className="text-start">
                    <h1 className="display-1">S.A.M.U.E.L.</h1>
                    <h2 className="h6 text-secondary-emphasis">
                        <strong>S</strong>uccessor <strong>A</strong>pplication for <strong>M</strong>anaging <strong>U</strong>ndergraduate <strong>E</strong>ducational <strong>L</strong>abors
                    </h2>
                </hgroup>
                <Button variant="secondary" size="lg" onClick={() => setShowNavigator(true)}>
                    <NavigatorMenuIcon />
                </Button>
            </Container>
        </>
    );
}

function Main({children}) {
    return (
        <Container className="d-flex flex-column p-3 gap-2">
            {children}
        </Container>
    );
}

export default function SAMUELHangar({user, session, children}) {
    return (
        <>
            <SessionExpirationModal session={session} />
            <Header user={user} />
            <Main>
                {children}
            </Main>
        </>
    );
}