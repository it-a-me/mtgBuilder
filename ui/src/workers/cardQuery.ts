import * as Comlink from "comlink"
import "./wasm_exec"

export default class CardQuery {
  async _init() {
    //@ts-expect-error untyped globalThis
    if (globalThis.GO_cardQuery != undefined) { return }

    //@ts-expect-error untyped wasm_exec
    const go = new Go();
    //@ts-expect-error untyped globalThis
    globalThis.GO_cardQuery = WebAssembly.instantiateStreaming(fetch("/wasmlib.wasm"), go.importObject).then((result) => {
      go.run(result.instance);
    });
    //@ts-expect-error untyped globalThis
    await globalThis.GO_cardQuery
  }

  async feedCards(cardsJson: string) {
    await this._init()
    console.log("initialized")
    // @ts-expect-error untyped value
    const res = self.GO_cardQuery.feedCards(cardsJson)
    if (res instanceof (Error)) {
      throw res
    }
    if (res == null) {
      return
    }
    throw "unreachable";
  }

  async queryCards(query: string): Promise<number[]> {
    // @ts-expect-error untyped value
    const res = globalThis.GO_cardQuery.queryCards(query)
    if (res instanceof (Error)) {
      throw res
    }
    if (Array.isArray(res)) {
      return res
    }
    throw "unreachable";
  }

  async parseQuery(query: string): Promise<object> {
    // @ts-expect-error untyped value
    const res = globalThis.GO_cardQuery.parseQuery(query)
    if (res instanceof (Error)) {
      throw res
    }
    if (typeof res == 'string') {
      return JSON.parse(res);
    }
    throw "unreachable";
  }

  async getCard(cardIndex: number): Promise<object> {
    // @ts-expect-error untyped value
    const res = globalThis.GO_cardQuery.getCard(cardIndex)
    if (res instanceof (Error)) {
      throw res
    }
    if (typeof res == 'string') {
      return JSON.parse(res);
    }
    throw new Error("unimplemented");
  }
}

const cardQuery = new CardQuery()
Comlink.expose(cardQuery)
