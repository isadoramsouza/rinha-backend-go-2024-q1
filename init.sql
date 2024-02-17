CREATE TABLE IF NOT EXISTS clientes (
    id SERIAL PRIMARY KEY NOT NULL,
    nome VARCHAR(50) NOT NULL,
    limite INTEGER NOT NULL,
    saldo INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS transacoes (
    id SERIAL PRIMARY KEY NOT NULL,
    tipo CHAR(1) NOT NULL,
    descricao VARCHAR(10) NOT NULL,
    valor INTEGER NOT NULL,
    cliente_id INTEGER NOT NULL,
    realizada_em TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cliente_id ON transacoes (cliente_id);

INSERT INTO clientes (nome, limite, saldo)
VALUES
    ('Isadora', 100000, 0),
    ('Maicon', 80000, 0),
    ('Matias', 1000000, 0),
    ('Bob', 10000000, 0),
    ('Tom', 500000, 0);

CREATE OR REPLACE FUNCTION atualizar_saldo()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE clientes 
    SET 
        saldo = CASE 
            WHEN NEW.tipo = 'd' THEN saldo - NEW.valor 
            ELSE saldo + NEW.valor 
        END
    WHERE 
        id = NEW.cliente_id 
        AND (NEW.tipo <> 'd' OR (saldo - NEW.valor) >= -limite);

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Débito excede o limite do cliente';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER atualizar_saldo_trigger
AFTER INSERT ON transacoes
FOR EACH ROW
EXECUTE FUNCTION atualizar_saldo();