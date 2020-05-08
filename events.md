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
