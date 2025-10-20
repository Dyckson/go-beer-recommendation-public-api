-- Seed data para popular a tabela beer_styles com 20 estilos de cerveja
-- 10 cervejas com temperaturas realistas + 10 cervejas com temperaturas variadas para testes

INSERT INTO beer_styles (name, temp_min, temp_max) VALUES 
    ('IPA', 7.0, 10.0),                    -- India Pale Ale: bem gelada
    ('Lager', 3.0, 6.0),                   -- Lager: muito gelada
    ('Stout', 10.0, 13.0),                 -- Stout: menos gelada, mais aromática
    ('Pilsner', 4.0, 7.0),                 -- Pilsner: bem gelada
    ('Wheat Beer', 6.0, 9.0),              -- Cerveja de trigo: moderadamente gelada
    ('Porter', 10.0, 13.0),                -- Porter: menos gelada para realçar sabores
    ('Pale Ale', 8.0, 11.0),               -- Pale Ale: moderadamente gelada
    ('Belgian Dubbel', 12.0, 15.0),        -- Belgian Dubbel: temperatura de adega
    ('Saison', 8.0, 11.0),                 -- Saison: moderadamente gelada
    ('Barleywine', 13.0, 16.0),            -- Barleywine: temperatura ambiente
    ('Arctic Lager', -5.0, 0.0),           -- Cerveja super gelada
    ('Frozen Ale', -10.0, -5.0),           -- Cerveja congelante
    ('Ice Beer', -2.0, 2.0),               -- Cerveja no limite do gelo
    ('Winter Porter', 0.0, 5.0),           -- Cerveja de inverno
    ('Hot Weather Lager', 20.0, 25.0),     -- Cerveja para clima quente
    ('Desert Ale', 30.0, 35.0),            -- Cerveja para clima muito quente
    ('Tropical IPA', 25.0, 30.0),          -- IPA tropical
    ('Extreme Cold', -15.0, -10.0),        -- Teste de temperatura extrema fria
    ('Warm Climate', 35.0, 40.0),          -- Teste de temperatura extrema quente
    ('Zero Point', -1.0, 1.0)              -- Cerveja no ponto zero
ON CONFLICT (name) DO NOTHING;
