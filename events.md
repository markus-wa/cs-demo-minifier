# Events

- [`fire`](#fire)
- [`hurt`](#hurt)
- [`kill`](#kill)
- [`flashed`](#flashed)
- [`jump`](#jump)
- [`footstep`](#footstep)
- [`chat_message`](#chat_message)
- [`swap_team`](#swap_team)
- [`disconnect`](#disconnect)
- [`round_started`](#round_started)
- [`round_ended`](#round_ended)
- [`smoke_started`](#smoke_started)
- [`smoke_expired`](#smoke_expired)
- [`decoy_started`](#decoy_started)
- [`decoy_expired`](#decoy_expired)
- [`fire_grenade_started`](#fire_grenade_started)
- [`fire_grenade_expired`](#fire_grenade_expired)
- [`flash_explosion`](#flash_explosion)
- [`flash_explosion`](#flash_explosion)

## Attributes

### `fire`

| attribute | type | description |
| --- | --- | --- |
| `entityId` | `numVal` | EntityID |

### `hurt`

| attribute | type | description |
| --- | --- | --- |
| `entityId` | `numVal` | EntityID |

### `kill`

| attribute | type | description |
| --- | --- | --- |
| `victim` | `numVal` | EntityID |
| `weapon` | `numVal` | see [`EquipmentType`](https://pkg.go.dev/github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common?tab=doc#EquipmentType) |
| `killer` | `numVal` | EntityID |
| `assister` | `numVal` | EntityID |

### `flashed`

| attribute | type | description |
| --- | --- | --- |
| `entityId` | `numVal` | EntityID |

### `jump`

| attribute | type | description |
| --- | --- | --- |
| `entityId` | `numVal` | EntityID |

### `footstep`

| attribute | type | description |
| --- | --- | --- |
| `entityId` | `numVal` | EntityID |

### `chat_message`

| attribute | type | description |
| --- | --- | --- |
| `sender` | `numVal` | EntityID |
| `text` | `strVal` | chat message |

### `swap_team`

| attribute | type | description |
| --- | --- | --- |
| `entityId` | `numVal` | EntityID |

### `disconnect`

| attribute | type | description |
| --- | --- | --- |
| `entityId` | `numVal` | EntityID |

### `round_started`

| attribute | type | description |
| --- | --- | --- |
| - | - | - |

### `round_ended`

| attribute | type | description |
| --- | --- | --- |
| `winner` | `numVal` | see [`Team`](https://pkg.go.dev/github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common?tab=doc#Team) |
| `reason` | `numVal` | see [`RoundEndReason`](https://pkg.go.dev/github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events?tab=doc#RoundEndReason) |

### `smoke_started`

| attribute | type | description |
| --- | --- | --- |
| `x` | `numVal` | The x-coordinate used in the CS:GO space |
| `y` | `numVal` | The y-coordinate used in the CS:GO space |
| `z` | `numVal` | The z-coordinate used in the CS:GO space |
| `throwerEntityId` | `numVal` | The entityId of the throwing player |

### `smoke_expired`

| attribute | type | description |
| --- | --- | --- |
| `x` | `numVal` | The x-coordinate used in the CS:GO space |
| `y` | `numVal` | The y-coordinate used in the CS:GO space |
| `z` | `numVal` | The z-coordinate used in the CS:GO space |
| `throwerEntityId` | `numVal` | The entityId of the throwing player |


### `decoy_started`

| attribute | type | description |
| --- | --- | --- |
| `x` | `numVal` | The x-coordinate used in the CS:GO space |
| `y` | `numVal` | The y-coordinate used in the CS:GO space |
| `z` | `numVal` | The z-coordinate used in the CS:GO space |
| `throwerEntityId` | `numVal` | The entityId of the throwing player |


### `decoy_expired`

| attribute | type | description |
| --- | --- | --- |
| `x` | `numVal` | The x-coordinate used in the CS:GO space |
| `y` | `numVal` | The y-coordinate used in the CS:GO space |
| `z` | `numVal` | The z-coordinate used in the CS:GO space |
| `throwerEntityId` | `numVal` | The entityId of the throwing player |

### `fire_grenade_started`

| attribute | type | description |
| --- | --- | --- |
| `x` | `numVal` | The x-coordinate used in the CS:GO space |
| `y` | `numVal` | The y-coordinate used in the CS:GO space |
| `z` | `numVal` | The z-coordinate used in the CS:GO space |

### `fire_grenade_expired`

| attribute | type | description |
| --- | --- | --- |
| `x` | `numVal` | The x-coordinate used in the CS:GO space |
| `y` | `numVal` | The y-coordinate used in the CS:GO space |
| `z` | `numVal` | The z-coordinate used in the CS:GO space |


### `fire_grenade_expired`

| attribute | type | description |
| --- | --- | --- |
| `x` | `numVal` | The x-coordinate used in the CS:GO space |
| `y` | `numVal` | The y-coordinate used in the CS:GO space |
| `z` | `numVal` | The z-coordinate used in the CS:GO space |


### `he_grenade_explosion`

| attribute | type | description |
| --- | --- | --- |
| `x` | `numVal` | The x-coordinate used in the CS:GO space |
| `y` | `numVal` | The y-coordinate used in the CS:GO space |
| `z` | `numVal` | The z-coordinate used in the CS:GO space |
| `throwerEntityId` | `numVal` | The entityId of the throwing player |

### `flash_explosion`

| attribute | type | description |
| --- | --- | --- |
| `x` | `numVal` | The x-coordinate used in the CS:GO space |
| `y` | `numVal` | The y-coordinate used in the CS:GO space |
| `z` | `numVal` | The z-coordinate used in the CS:GO space |
| `throwerEntityId` | `numVal` | The entityId of the throwing player |

