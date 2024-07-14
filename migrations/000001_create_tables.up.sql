CREATE TABLE dishes (
    id UUID PRIMARY KEY,
    kitchen_id UUID REFERENCES kitchens(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(50),
    ingredients TEXT[],
    allergens TEXT[],
    nutrition_info JSONB,
    dietary_info TEXT[],
    available BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    kitchen_id UUID REFERENCES kitchens(id),
    items JSONB NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    delivery_address TEXT NOT NULL,
    delivery_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE reviews (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),
    user_id UUID REFERENCES users(id),
    kitchen_id UUID REFERENCES kitchens(id),
    rating DECIMAL(2, 1) NOT NULL,
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE payments (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    transaction_id VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE working_hours (
    kitchen_id UUID REFERENCES kitchens(id),
    day_of_week INTEGER NOT NULL,
    open_time TIME NOT NULL,
    close_time TIME NOT NULL,
    PRIMARY KEY (kitchen_id, day_of_week)
);

CREATE TABLE delivery_routes (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),
    start_address TEXT NOT NULL,
    end_address TEXT NOT NULL,
    distance DECIMAL(10, 2) NOT NULL,
    duration INTEGER NOT NULL,
    route_polyline TEXT,
    waypoints JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, 
    deleted_at TIMESTAMP
);
