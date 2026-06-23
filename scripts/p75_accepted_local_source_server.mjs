import http from 'node:http'
import { writeFileSync } from 'node:fs'

const portFile = process.argv[2]
const requestLogFile = process.argv[3]

if (!portFile || !requestLogFile) {
  console.error('usage: node scripts/p75_accepted_local_source_server.mjs <port-file> <request-log-file>')
  process.exit(2)
}

const requests = []

function logRequest(req, body) {
  const url = new URL(req.url || '/', 'http://127.0.0.1')
  const form = parseForm(body)
  requests.push({
    method: req.method,
    path: url.pathname,
    query: Object.fromEntries(url.searchParams.entries()),
    form,
  })
  writeFileSync(requestLogFile, `${JSON.stringify(requests, null, 2)}\n`, 'utf8')
  return url
}

function parseForm(body) {
  const text = body.toString('utf8')
  if (!text.includes('=')) return {}
  return Object.fromEntries(new URLSearchParams(text).entries())
}

function sendJSON(res, value) {
  res.writeHead(200, { 'Content-Type': 'application/json; charset=utf-8' })
  res.end(`${JSON.stringify(value)}\n`)
}

const server = http.createServer((req, res) => {
  const chunks = []
  req.on('data', (chunk) => chunks.push(chunk))
  req.on('end', () => handleRequest(req, res, Buffer.concat(chunks)))
})

function handleRequest(req, res, body) {
  const url = logRequest(req, body)
  if (url.pathname === '/market') {
    if (url.searchParams.get('symbol') !== '159915') {
      res.writeHead(400, { 'Content-Type': 'application/json' })
      res.end('{"error":"unexpected symbol"}')
      return
    }
    sendJSON(res, {
      close_price: 2.413,
      turnover_rate: 1.4,
      pe_percentile: 42,
      pb_percentile: 47,
      volume_percentile: 51,
      volatility_percentile: 33,
      liquidity_state: 'normal',
      sentiment_state: 'neutral',
      source_name: 'accepted_local_market',
      source_level: 'B',
      source_type: 'market_price',
      trade_date: '2026-06-19',
      captured_at: '2026-06-19T08:00:00Z',
      metadata: {
        p34_source_health: {
          symbol_profile: { freshness: 'fresh', data_date: '2026-06-19', affected_symbols: ['159915'], source_level: 'A', source_name: 'accepted_local_registry', source_type: 'symbol_profile' },
          fund_profile: { freshness: 'fresh', data_date: '2026-06-19', affected_symbols: ['159915'], source_level: 'B', source_name: 'accepted_local_fund_profile', source_type: 'fund_profile' },
          tracked_index: { freshness: 'fresh', data_date: '2026-06-19', affected_symbols: ['399006'], source_level: 'A', source_name: 'accepted_local_index_profile', source_type: 'index_profile' },
          market_price: { freshness: 'fresh', data_date: '2026-06-19', affected_symbols: ['159915'], source_level: 'B', source_name: 'accepted_local_market', source_type: 'market_price' },
          valuation_percentiles: { freshness: 'fresh', data_date: '2026-06-19', affected_symbols: ['399006'], source_level: 'A', source_name: 'accepted_local_index_valuation', source_type: 'index_valuation' },
          liquidity: { freshness: 'fresh', data_date: '2026-06-19', affected_symbols: ['159915'], source_level: 'B', source_name: 'accepted_local_liquidity', source_type: 'liquidity' },
          sentiment_proxy: { freshness: 'fresh', data_date: '2026-06-19', affected_symbols: ['159915'], source_level: 'C', source_name: 'accepted_local_sentiment', source_type: 'sentiment_proxy' },
          rag_index: { freshness: 'fresh', data_date: '2026-06-19', affected_symbols: ['159915'], source_level: 'local_index', source_name: 'veclite', source_type: 'rag_index' },
        },
        p34_data_categories: ['symbol_profile', 'fund_profile', 'tracked_index', 'market_price', 'valuation_percentiles', 'liquidity', 'sentiment_proxy', 'rag_index'],
      },
    })
    return
  }
  if (url.pathname === '/new/hisAnnouncement/query') {
    sendJSON(res, {
      totalAnnouncement: 1,
      totalRecordNum: 1,
      announcements: [{
        secCode: '159915',
        secName: '创业板ETF',
        orgId: 'org',
        announcementId: 'ann-159915',
        announcementTitle: '创业板ETF 公告',
        announcementTime: 1781827200000,
        adjunctUrl: '/159915.pdf',
        adjunctType: 'PDF',
        announcementType: '基金公告',
      }],
      hasMore: false,
      totalpages: 1,
    })
    return
  }
  if (url.pathname === '/api/disc/announcement/searchQuery') {
    sendJSON(res, {
      recordCount: 1,
      data: [{
        secCode: '159915',
        secName: '创业板ETF',
        announList: [{
          id: 'szse-159915',
          title: '创业板ETF 公告',
          attachPath: '/159915-szse.pdf',
          attachFormat: 'PDF',
          attachSize: 1,
          annId: 'ann-159915-szse',
          bigCategoryName: '定期报告',
          publishTime: '2026-06-19 09:00:00',
        }],
      }],
    })
    return
  }
  res.writeHead(404, { 'Content-Type': 'application/json' })
  res.end('{"error":"not found"}')
}

server.listen(0, '127.0.0.1', () => {
  const address = server.address()
  writeFileSync(portFile, `${address.port}\n`, 'utf8')
})
