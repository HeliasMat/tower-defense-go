tower-defense/
├── cmd/
│   └── main.go
├── internal/
│   ├── core/
│   │   ├── engine.go
│   │   ├── game_state.go
│   │   ├── config.go
│   │   └── constants.go
│   ├── entities/
│   │   ├── entity.go
│   │   ├── tower.go
│   │   ├── enemy.go
│   │   ├── projectile.go
│   │   └── effects/
│   │       ├── particle.go
│   │       └── animation.go
│   ├── systems/
│   │   ├── movement_system.go
│   │   ├── collision_system.go
│   │   ├── targeting_system.go
│   │   ├── spawning_system.go
│   │   └── cleanup_system.go
│   ├── utils/
│   │   ├── math.go
│   │   ├── spatial_hash.go
│   │   └── object_pool.go
│   ├── rendering/
│   │   ├── renderer.go
│   │   ├── sprite_manager.go
│   │   └── shaders/
│   │       ├── vertex_shader.glsl
│   │       └── fragment_shader.glsl
│   ├── audio/
│   │   ├── audio_manager.go
│   │   └── sound_effects.go
│   └── ui/
│       ├── hud.go
│       └── menu.go
├── pkg/
│   └── // Placeholder for shared packages
├── tests/
│   ├── unit/
│   │   ├── core/
│   │   │   ├── engine_test.go
│   │   │   ├── game_state_test.go
│   │   │   ├── config_test.go
│   │   │   └── constants_test.go
│   │   ├── entities/
│   │   │   ├── entity_test.go
│   │   │   ├── tower_test.go
│   │   │   ├── enemy_test.go
│   │   │   ├── projectile_test.go
│   │   │   └── effects/
│   │   │       ├── particle_test.go
│   │   │       └── animation_test.go
│   │   ├── systems/
│   │   │   ├── movement_system_test.go
│   │   │   ├── collision_system_test.go
│   │   │   ├── targeting_system_test.go
│   │   │   ├── spawning_system_test.go
│   │   │   └── cleanup_system_test.go
│   │   ├── utils/
│   │   │   ├── math_test.go
│   │   │   ├── spatial_hash_test.go
│   │   │   └── object_pool_test.go
│   │   ├── rendering/
│   │   │   ├── renderer_test.go
│   │   │   ├── sprite_manager_test.go
│   │   │   └── shaders/
│   │   │       ├── vertex_shader_test.go
│   │   │       └── fragment_shader_test.go
│   │   ├── audio/
│   │   │   ├── audio_manager_test.go
│   │   │   └── sound_effects_test.go
│   │   └── ui/
│   │       ├── hud_test.go
│   │       └── menu_test.go
│   ├── integration/
│   │   ├── game_loop_test.go
│   │   ├── performance_test.go
│   │   ├── core_systems_integration_test.go
│   │   ├── entity_interaction_test.go
│   │   └── rendering_audio_integration_test.go
│   ├── benchmarks/
│   │   ├── core/
│   │   │   ├── engine_benchmark_test.go
│   │   │   ├── game_state_benchmark_test.go
│   │   ├── entities/
│   │   │   ├── entity_benchmark_test.go
│   │   │   ├── tower_benchmark_test.go
│   │   │   ├── enemy_benchmark_test.go
│   │   │   ├── projectile_benchmark_test.go
│   │   ├── systems/
│   │   │   ├── movement_system_benchmark_test.go
│   │   │   ├── collision_system_benchmark_test.go
│   │   │   ├── targeting_system_benchmark_test.go
│   │   │   ├── spawning_system_benchmark_test.go
│   │   └── rendering/
│   │       ├── renderer_benchmark_test.go
│   │       └── sprite_manager_benchmark_test.go
├── assets/
│   ├── textures/
│   ├── models/
│   ├── audio/
│   └── maps/
├── configs/
│   ├── game_config.yaml
│   └── enemy_waves.json
├── docs/
│   ├── architecture.md
│   ├── optimization_strategies.md
│   └── api/
├── tools/
│   ├── map_editor/
│   └── profiling/
├── go.mod
├── go.sum
└── README.md
