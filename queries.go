package main

// Subscription counts, regardless of current member status - ie some may be deceased or retired
const QUERY_SUBSCRIPTION_COUNTS = `SELECT
  s.name   AS Subscription,
  COUNT(*) AS Count
FROM fn_m_subscription ms
  LEFT JOIN fn_subscription s ON ms.fn_subscription_id = s.id
GROUP BY s.id
ORDER BY Subscription ASC`


// Subscription count for currently active members
const QUERY_ACTIVE_SUBSCRIPTION_COUNTS = `SELECT 
    s.name AS Subscription, COUNT(*) AS Count 
FROM
    fn_m_subscription ms
        LEFT JOIN
    fn_subscription s ON ms.fn_subscription_id = s.id
WHERE
    ms.member_id IN (SELECT 
            member_id
        FROM
            ms_m_status
        WHERE
            ms_status_id = 1 AND active = 1
                AND current = 1)
GROUP BY s.id
ORDER BY Subscription ASC`


const QUERY_MEMBER_ID = `SELECT id FROM member`

const QUERY_MEMBER_STATUS_HISTORY = `SELECT
	DATE_FORMAT(ms.created_at, "%Y-%m-%d") as Date,
    ms.ms_status_id AS StatusID,
    s.name as Name
FROM
    ms_m_status ms
        LEFT JOIN
    ms_status s ON ms.ms_status_id = s.id
WHERE
    member_id = ?
ORDER BY ms.created_at ASC`

const QUERY_MEMBER_TITLE_HISTORY = `SELECT
	mt.granted_on as Date,
    mt.ms_title_id AS TitleID,
    t.name as Name
FROM
    ms_m_title mt
        LEFT JOIN
    ms_title t ON mt.ms_title_id = t.id
WHERE
    member_id = ?
ORDER BY mt.granted_on ASC`

const QUERY_TITLES = `SELECT id, name FROM ms_title;`

const QUERY_CURRENTLY_LAPSED_MEMBERS_TITLE_YEAR = `SELECT 
    t.name AS Title,
    ms.member_id AS MemberID,
    DATE_FORMAT(ms.updated_at, '%Y') AS LapsedOn
FROM
    ms_m_status ms
        LEFT JOIN
    ms_m_title mt ON ms.member_id = mt.member_id
        LEFT JOIN
    ms_title t ON mt.ms_title_id = t.id
WHERE
    ms.ms_status_id = 10004
        AND ms.comment != 'Initial data import'
        AND t.name != 'Admin Member'
        AND ms.current = 1
		AND mt.current = 1`

const QUERY_CURRENTLY_LAPSED_MEMBERS_COUNT_TITLE_YEAR = `SELECT 
	DATE_FORMAT(ms.updated_at, '%Y') AS LapsedYear,
    t.name AS Title,
    count(ms.member_id) AS Count
FROM
    ms_m_status ms
        LEFT JOIN
    ms_m_title mt ON ms.member_id = mt.member_id
        LEFT JOIN
    ms_title t ON mt.ms_title_id = t.id
WHERE
    ms.ms_status_id = 10004
        AND ms.comment != 'Initial data import'
        AND t.name != 'Admin Member'
        AND ms.current = 1
		AND mt.current = 1
GROUP BY LapsedYear,Title`