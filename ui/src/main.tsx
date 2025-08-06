import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App'
import CardQuery from './cardQuery'

const c = new CardQuery();
async function fetchCards() {
  const cardsJSON = await fetch("/cards.json").then(r => r.text())
  c.feedCards(cardsJSON)
}
await fetchCards()

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App cardQuery={c} />
  </StrictMode>,
)
