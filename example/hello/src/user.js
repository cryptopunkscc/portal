import {rpc, log} from 'portal';

const ASTRAL_TARGET = "localnode";

export async function opUser() {
    try {
        const response = await rpc.target(ASTRAL_TARGET)
            .call("user.info")
            .request({"out": "json"});


        const { UserAlias, NodeAlias } = response?.Object ?? {};

        return {
            hasUser: Boolean(UserAlias),
            userAlias: UserAlias ?? 'Unknown',
            nodeAlias: NodeAlias ?? 'Unknown'
        };
    } catch (e) {
        log(e);
        return {
            hasUser: false,
            userAlias: null,
            nodeAlias: null
        };
    }
}

export async function opSwarmStatus() {
    try {
        const res = await rpc.target(ASTRAL_TARGET)
            .call("user.swarm_status")
            .collect({"out": "json"})

        var objects = res.map(el => {
            return {
                identity: el.Object.Identity,
                alias:  el.Object.Alias,
                linked:  el.Object.Linked,
                expiresAt:  el.Object.Contract?.ExpiresAt
            }
        })

        return {
            success: true,
            members: objects,
        }
    } catch (e) {
        console.log(e)
        return {
            success: false,
            members: []
        };
    }
}
