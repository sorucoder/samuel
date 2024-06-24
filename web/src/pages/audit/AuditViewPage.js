import React, { useState } from 'react';
import { useLoaderData } from 'react-router-dom';
import { DateTime } from 'luxon'; 
import { sessionGET, errorFromResponse } from '../../api.js';

import SAMUELBoard from '../../components/control/SAMUELBoard.js';
import SAMUELHangar from '../../components/page/SAMUELHangar.js';
import SAMUELDatePicker from '../../components/control/SAMUELDatePicker.js';

const ITEMS_PER_PAGE = 30;

async function fetchAudits({date = DateTime.now(), page = 0} = {}) {
    const apiResponse = await sessionGET('audit/view', {date: date.toISODate(), page: page, count: ITEMS_PER_PAGE});
    if (apiResponse.ok) {
        const {user, session, batch} = await apiResponse.json();
        return {
            user: {
                role: user.role
            },
            session: {
                expiresOn: DateTime.fromISO(session.expiresOn)
            },
            date: date,
            pagination: {
                number: batch.number,
                count: batch.count
            },
            audits: {
                columns: ['Description', 'User', 'Timestamp'],
                rows: batch.audits.map(({description, userUUID, timestamp}) => [description, userUUID, timestamp])
            }
        };
    } else {
        switch (apiResponse.status) {
            case 401:
                throw await errorFromResponse(apiResponse, 'Your Session Has Expired', { path: '/login', name: 'Login Page' });
            case 403:
                throw await errorFromResponse(apiResponse, 'You Are Not Authorized', { path: '/dashboard', name: 'Dashboard' });
            default:
                throw await errorFromResponse(apiResponse, 'Failed To View Audit', { path: '/dashboard', name: 'Dashboard' });
        }
    }
}

export async function auditViewLoader() {
    return fetchAudits();
}

export default function AuditViewPage() {
    const loaderData = useLoaderData();
    
    const [loading, setLoading] = useState(false);
    const [session, setSession] = useState(loaderData.session);
    const [audits, setAudits] = useState(loaderData.audits);
    const [date, setDate] = useState(loaderData.date);
    const [pagination, setPagination] = useState(loaderData.pagination);
    const onDateChange = async (newDate) => {
        setLoading(true);

        const {session, pagination, audits} = await fetchAudits({date: newDate});
        setSession(session);
        setDate(newDate);
        setPagination(pagination);
        setAudits(audits);

        setLoading(false);
    };
    const onPageChange = async (newPage) => {
        setLoading(true);
        
        const {session, pagination, audits} = await fetchAudits({date: date, page: newPage})
        setSession(session);
        setPagination(pagination);
        setAudits(audits);

        setLoading(false);
    }
    
    return (
        <SAMUELHangar user={loaderData.user} session={session}>
            <h3>Audit</h3>
            <SAMUELDatePicker className="w-25" name="date" value={date} onDateChange={onDateChange} />
            <SAMUELBoard loading={loading} noDataMessage="There are no audits for this day." data={audits} pagination={pagination} onPageChange={onPageChange} />
        </SAMUELHangar>
    );
}