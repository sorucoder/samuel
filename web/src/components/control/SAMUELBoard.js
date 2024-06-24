import React from 'react';

import Pagination from 'react-bootstrap/Pagination';
import Spinner from 'react-bootstrap/Spinner';
import Table from 'react-bootstrap/Table';

const MAXIMUM_PAGINATION_ITEMS = 3;

function TableHeader({columns}) {
    return (
        <thead>
            <tr>
                {columns.map((column, columnIndex) => (
                    <th key={columnIndex} className="text-bg-primary">
                        {column}
                    </th>
                ))}
            </tr>
        </thead>
    );
}

function TableBody({data, noDataMessage}) {
    if (data.rows.length === 0) {
        return (
            <tbody>
                <tr>
                    <td className="text-center" colspan={data.columns.length}>
                        <em>{noDataMessage ?? 'There are no entries.'}</em>
                    </td>
                </tr>
            </tbody>
        );
    } else {
        return (
            <tbody>
                {data.rows.map((row, rowIndex) => (
                    <tr key={rowIndex}>
                        {row.map((cell, columnIndex) => (
                            <td key={columnIndex}>
                                {cell}
                            </td>
                        ))}
                    </tr>
                ))}
            </tbody>
        );
    }
}

function TablePagination({pagination, onPageChange}) {
    if (pagination.count === 0) {
        return null;
    } else {
        const onFirst = pagination.number === 0;
        const onLast = pagination.number === pagination.count - 1;
        const items = new Array(9);
        items.push(
            <Pagination.First key="first" disabled={onFirst} onClick={() => onPageChange(0)} />
        )
        if (pagination.count > MAXIMUM_PAGINATION_ITEMS) {
            items.push(
                <Pagination.Prev key="previous" disabled={onFirst} onClick={() => onPageChange(Math.max(0, pagination.number - 1))} />
            );
            
            if (pagination.number < MAXIMUM_PAGINATION_ITEMS) {
                for (let number = 0; number < MAXIMUM_PAGINATION_ITEMS; number++) {
                    items.push(
                        <Pagination.Item key={number} active={number === pagination.number} onClick={() => onPageChange(number)}>{number + 1}</Pagination.Item>
                    );
                }
                items.push(
                    <Pagination.Ellipsis key="next-ellipsis" />
                );
            } else if (pagination.number >= pagination.count - MAXIMUM_PAGINATION_ITEMS) {
                items.push(
                    <Pagination.Ellipsis key="previous-ellipsis" />
                );
                for (let number = pagination.count - MAXIMUM_PAGINATION_ITEMS; number < pagination.count; number++) {
                    items.push(
                        <Pagination.Item key={number} active={number === pagination.number} onClick={() => onPageChange(number)}>{number + 1}</Pagination.Item>
                    );
                }
            } else {
                items.push(
                    <Pagination.Ellipsis key="previous-ellipsis" />
                );
                for (let number = pagination.number - 1; number <= pagination.number + 1; number++) {
                    items.push(
                        <Pagination.Item key={number} active={number === pagination.number} onClick={() => onPageChange(number)}>{number + 1}</Pagination.Item>
                    );
                }
                items.push(
                    <Pagination.Ellipsis key="next-ellipsis" />
                );
            }

            items.push(
                <Pagination.Next key="next" disabled={onLast} onClick={() => onPageChange(Math.min(pagination.count - 1, pagination.number + 1))} />
            );
        } else {
            for (let index = 0; index < pagination.count; index++) {
                items.push(
                    <Pagination.Item key={index} active={index === pagination.number} onClick={() => onPageChange(index)}>{index + 1}</Pagination.Item>
                );
            }
        }
        items.push(
            <Pagination.Last key="last" disabled={onLast} onClick={() => onPageChange(pagination.count - 1)} />
        );

        return (
            <Pagination className="align-self-center">
                {items}
            </Pagination>
        );
    }
}

export default function SAMUELManagementTable({loading, noDataMessage, data, pagination, onPageChange}) {
    if (loading) {
        return (
            <Spinner animation="border" variant="primary" className="align-self-center" role="status">
                <span className="visually-hidden">Loading...</span>
            </Spinner>
        );
    } else {
        return (
            <>
                <Table hover striped>
                    <TableHeader columns={data.columns} />
                    <TableBody data={data} noDataMessage={noDataMessage} />
                </Table>
                <TablePagination pagination={pagination} onPageChange={onPageChange} />
            </>
        );
    }
}