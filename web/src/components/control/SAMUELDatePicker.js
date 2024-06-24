import React from 'react';
import { DateTime, Duration } from 'luxon';

import Button from 'react-bootstrap/Button';
import FormControl from 'react-bootstrap/FormControl';

import InputGroup from 'react-bootstrap/InputGroup';

import { ChevronLeft as PreviousDayIcon, ChevronRight as NextDayIcon } from 'react-bootstrap-icons';

const ONE_DAY = Duration.fromObject({days: 1});

export default function SAMUELDatePicker({name, label, value, onDateChange, className}) {
    
    return (
        <InputGroup className={className}>
            <Button variant="primary" onClick={() => onDateChange(value.minus(ONE_DAY))}>
                <PreviousDayIcon />
            </Button>
            <FormControl type="date" name={name} placeholder={label} value={value.toISODate()} onChange={(event) => onDateChange(DateTime.fromISO(event.target.value))} />
            <Button variant="primary" onClick={() => onDateChange(value.plus(ONE_DAY))}>
                <NextDayIcon />
            </Button>
        </InputGroup>
    );
}